package grpcdapter

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/beabys/wms/auth-service/internal/application/auth/repository"
	"github.com/beabys/wms/auth-service/internal/application/auth/usecase"
	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
	repomocks "github.com/beabys/wms/auth-service/mocks/application/auth/repository"
	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
	"github.com/beabys/wms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type serverDeps struct {
	userRepo   *repomocks.UserRepository
	authRepo   *repomocks.AuthRepository
	inviteRepo *repomocks.InviteRepository
}

func setupGRPCServer(t *testing.T, deps serverDeps) (*grpc.ClientConn, *usecase.Service) {
	t.Helper()
	log := logger.NewSlogLogger(slog.LevelDebug)

	privPEM, pubPEM, err := usecase.GenerateKeyPair()
	require.NoError(t, err)

	jwtSvc, err := usecase.NewService(privPEM, pubPEM, 15*time.Minute, 7*24*time.Hour)
	require.NoError(t, err)

	loginUC := usecase.NewLoginUseCase(log, deps.userRepo, deps.authRepo, jwtSvc)
	userUC := usecase.NewUserUseCase(log, deps.userRepo)
	inviteUC := usecase.NewInviteUseCase(log, deps.inviteRepo, deps.userRepo)

	// Create in-process gRPC server
	lis := bufconn.Listen(1024 * 1024)

	s := grpc.NewServer()
	authv1.RegisterAuthServiceServer(s, NewAuthServer(loginUC, inviteUC))
	authv1.RegisterUserServiceServer(s, NewUserServer(userUC))

	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// Dial in-process
	conn, err := grpc.Dial("bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		conn.Close()
		s.Stop()
	})

	return conn, jwtSvc
}

func TestAuthServer_Login(t *testing.T) {
	t.Run("login with valid credentials returns tokens", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)
		userClient := authv1.NewUserServiceClient(conn)

		user, _ := model.NewUser("logintest@example.com", "password123", "Login Test", "admin", nil)

		// CreateUser expectations
		userRepo.On("GetByEmail", mock.Anything, "logintest@example.com").Return(nil, nil).Once()
		userRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Email == "logintest@example.com"
		})).Return(nil).Once()

		_, err := userClient.CreateUser(context.Background(), &authv1.CreateUserRequest{
			Email:    "logintest@example.com",
			Password: "password123",
			Name:     "Login Test",
			Role:     "admin",
		})
		require.NoError(t, err)

		// Login expectations
		userRepo.On("GetByEmail", mock.Anything, "logintest@example.com").Return(user, nil).Once()
		authRepo.On("StoreRefreshToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		resp, err := client.Login(context.Background(), &authv1.LoginRequest{
			Email:    "logintest@example.com",
			Password: "password123",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, resp.GetAccessToken())
		assert.NotEmpty(t, resp.GetRefreshToken())
		assert.Greater(t, resp.GetExpiresIn(), int64(0))
	})

	t.Run("login with wrong password returns error", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)

		user, _ := model.NewUser("logintest@example.com", "password123", "Login Test", "admin", nil)

		userRepo.On("GetByEmail", mock.Anything, "logintest@example.com").Return(user, nil).Once()

		_, err := client.Login(context.Background(), &authv1.LoginRequest{
			Email:    "logintest@example.com",
			Password: "wrongpassword",
		})
		require.Error(t, err)
	})
}

func TestAuthServer_ValidateToken(t *testing.T) {
	t.Run("valid token returns valid response", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, jwtSvc := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)

		user := &model.User{ID: "validate-user-1", Role: "admin", Email: "v@example.com", Name: "V"}
		token, err := jwtSvc.GenerateAccessToken(user)
		require.NoError(t, err)

		resp, err := client.ValidateToken(context.Background(), &authv1.ValidateTokenRequest{
			Token: token,
		})
		require.NoError(t, err)
		assert.True(t, resp.GetValid())
		assert.Equal(t, "validate-user-1", resp.GetUserId())
		assert.Equal(t, "admin", resp.GetRole())
	})

	t.Run("invalid token returns valid=false", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)

		resp, err := client.ValidateToken(context.Background(), &authv1.ValidateTokenRequest{
			Token: "invalid.jwt.token",
		})
		require.NoError(t, err)
		assert.False(t, resp.GetValid())
	})
}

func TestAuthServer_Logout(t *testing.T) {
	t.Run("logout with valid refresh token revokes it", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)
		userClient := authv1.NewUserServiceClient(conn)

		user, _ := model.NewUser("logouttest@example.com", "password123", "Logout Test", "admin", nil)

		// CreateUser expectations
		userRepo.On("GetByEmail", mock.Anything, "logouttest@example.com").Return(nil, nil).Once()
		userRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Email == "logouttest@example.com"
		})).Return(nil).Once()

		_, err := userClient.CreateUser(context.Background(), &authv1.CreateUserRequest{
			Email:    "logouttest@example.com",
			Password: "password123",
			Name:     "Logout Test",
			Role:     "admin",
		})
		require.NoError(t, err)

		// Login expectations
		userRepo.On("GetByEmail", mock.Anything, "logouttest@example.com").Return(user, nil).Once()
		authRepo.On("StoreRefreshToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		loginResp, err := client.Login(context.Background(), &authv1.LoginRequest{
			Email:    "logouttest@example.com",
			Password: "password123",
		})
		require.NoError(t, err)
		require.NotEmpty(t, loginResp.GetRefreshToken())

		// Logout expectations
		authRepo.On("RevokeRefreshToken", mock.Anything, mock.Anything).Return(nil).Once()

		logoutResp, err := client.Logout(context.Background(), &authv1.LogoutRequest{
			RefreshToken: loginResp.GetRefreshToken(),
		})
		require.NoError(t, err)
		assert.True(t, logoutResp.GetSuccess())

		authRepo.AssertCalled(t, "RevokeRefreshToken", mock.Anything, mock.Anything)
	})

	t.Run("logout with empty token returns success", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)

		resp, err := client.Logout(context.Background(), &authv1.LogoutRequest{
			RefreshToken: "",
		})
		require.NoError(t, err)
		assert.True(t, resp.GetSuccess())
	})
}

func TestAuthServer_ValidateInvite(t *testing.T) {
	t.Run("valid invite token returns valid with email", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)

		// InviteUser expectations
		userRepo.On("GetByEmail", mock.Anything, "invited@example.com").Return(nil, nil).Once()
		inviteRepo.On("ExistsPendingByEmail", mock.Anything, "invited@example.com").Return(false, nil).Once()
		inviteRepo.On("Create", mock.Anything, mock.MatchedBy(func(inv *model.InviteToken) bool {
			return inv.Email == "invited@example.com"
		})).Return(nil).Once()

		inviteResp, err := client.InviteUser(context.Background(), &authv1.InviteUserRequest{
			Email:     "invited@example.com",
			InvitedBy: "admin-123",
		})
		require.NoError(t, err)
		require.NotEmpty(t, inviteResp.GetToken())

		// ValidateInvite expectations
		invite := &model.InviteToken{
			ID:        "inv-1",
			Email:     "invited@example.com",
			Token:     inviteResp.GetToken(),
			InvitedBy: "admin-123",
			Status:    model.StatusPending,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now(),
		}
		inviteRepo.On("GetByToken", mock.Anything, inviteResp.GetToken()).Return(invite, nil).Once()
		inviteRepo.On("MarkUsed", mock.Anything, inviteResp.GetToken()).Return(nil).Once()

		resp, err := client.ValidateInvite(context.Background(), &authv1.ValidateInviteRequest{
			Token: inviteResp.GetToken(),
		})
		require.NoError(t, err)
		assert.True(t, resp.GetValid())
		assert.Equal(t, "invited@example.com", resp.GetEmail())
		assert.Equal(t, "admin-123", resp.GetInvitedBy())
	})

	t.Run("expired token returns error", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)

		inviteRepo.On("GetByToken", mock.Anything, "nonexistent-token").Return(nil, nil).Once()

		_, err := client.ValidateInvite(context.Background(), &authv1.ValidateInviteRequest{
			Token: "nonexistent-token",
		})
		require.Error(t, err)
	})

	t.Run("empty token returns error", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)

		_, err := client.ValidateInvite(context.Background(), &authv1.ValidateInviteRequest{
			Token: "",
		})
		require.Error(t, err)
	})
}

func TestAuthServer_ListInvites(t *testing.T) {
	t.Run("list all invites returns paginated results", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)

		// For each InviteUser call (3), set up expectations
		for i := 0; i < 3; i++ {
			email := fmt.Sprintf("list%d@example.com", i)
			userRepo.On("GetByEmail", mock.Anything, email).Return(nil, nil).Once()
			inviteRepo.On("ExistsPendingByEmail", mock.Anything, email).Return(false, nil).Once()
			inviteRepo.On("Create", mock.Anything, mock.MatchedBy(func(inv *model.InviteToken) bool {
				return inv.Email == email
			})).Return(nil).Once()
		}

		for i := 0; i < 3; i++ {
			_, err := client.InviteUser(context.Background(), &authv1.InviteUserRequest{
				Email:     fmt.Sprintf("list%d@example.com", i),
				InvitedBy: "admin-123",
			})
			require.NoError(t, err)
		}

		// List expects
		inviteRepo.On("List", mock.Anything, repository.InviteFilter{Page: 1, PageSize: 10}).Return(&repository.InviteListResult{
			Invites: []*model.InviteToken{
				{ID: "1", Email: "list0@example.com", Token: "tok-0", InvitedBy: "admin-123", Status: model.StatusPending, ExpiresAt: time.Now().Add(24 * time.Hour), CreatedAt: time.Now()},
				{ID: "2", Email: "list1@example.com", Token: "tok-1", InvitedBy: "admin-123", Status: model.StatusPending, ExpiresAt: time.Now().Add(24 * time.Hour), CreatedAt: time.Now()},
				{ID: "3", Email: "list2@example.com", Token: "tok-2", InvitedBy: "admin-123", Status: model.StatusPending, ExpiresAt: time.Now().Add(24 * time.Hour), CreatedAt: time.Now()},
			},
			TotalItems: 3,
		}, nil).Once()

		resp, err := client.ListInvites(context.Background(), &authv1.ListInvitesRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(resp.GetInvites()), 3)
		assert.Equal(t, int32(1), resp.GetPagination().GetPage())
		assert.Equal(t, int32(10), resp.GetPagination().GetPageSize())
		assert.GreaterOrEqual(t, resp.GetPagination().GetTotal(), int32(3))
	})

	t.Run("filter by status used", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)

		status := "used"
		inviteRepo.On("List", mock.Anything, repository.InviteFilter{
			Status:   &status,
			Page:     1,
			PageSize: 10,
		}).Return(&repository.InviteListResult{
			Invites:    []*model.InviteToken{},
			TotalItems: 0,
		}, nil).Once()

		resp, err := client.ListInvites(context.Background(), &authv1.ListInvitesRequest{
			Page:     1,
			PageSize: 10,
			Status:   &status,
		})
		require.NoError(t, err)
		assert.Equal(t, int32(0), resp.GetPagination().GetTotal())
	})

	t.Run("filter by non-expired invites", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)

		expired := false
		inviteRepo.On("List", mock.Anything, repository.InviteFilter{
			Expired:  &expired,
			Page:     1,
			PageSize: 10,
		}).Return(&repository.InviteListResult{
			Invites: []*model.InviteToken{
				{ID: "1", Email: "list0@example.com", Token: "tok-0", InvitedBy: "admin-123", Status: model.StatusPending, ExpiresAt: time.Now().Add(24 * time.Hour), CreatedAt: time.Now()},
				{ID: "2", Email: "list1@example.com", Token: "tok-1", InvitedBy: "admin-123", Status: model.StatusPending, ExpiresAt: time.Now().Add(24 * time.Hour), CreatedAt: time.Now()},
				{ID: "3", Email: "list2@example.com", Token: "tok-2", InvitedBy: "admin-123", Status: model.StatusPending, ExpiresAt: time.Now().Add(24 * time.Hour), CreatedAt: time.Now()},
			},
			TotalItems: 3,
		}, nil).Once()

		resp, err := client.ListInvites(context.Background(), &authv1.ListInvitesRequest{
			Page:    1,
			PageSize: 10,
			Expired: &expired,
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(resp.GetInvites()), 3)
	})

	t.Run("filter by created after timestamp", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)

		after := time.Now().Add(-1 * time.Hour).Unix()
		createdAfter := time.Unix(after, 0)
		inviteRepo.On("List", mock.Anything, repository.InviteFilter{
			CreatedAfter: &createdAfter,
			Page:         1,
			PageSize:     10,
		}).Return(&repository.InviteListResult{
			Invites: []*model.InviteToken{
				{ID: "1", Email: "list0@example.com", Token: "tok-0", InvitedBy: "admin-123", Status: model.StatusPending, ExpiresAt: time.Now().Add(24 * time.Hour), CreatedAt: time.Now()},
				{ID: "2", Email: "list1@example.com", Token: "tok-1", InvitedBy: "admin-123", Status: model.StatusPending, ExpiresAt: time.Now().Add(24 * time.Hour), CreatedAt: time.Now()},
				{ID: "3", Email: "list2@example.com", Token: "tok-2", InvitedBy: "admin-123", Status: model.StatusPending, ExpiresAt: time.Now().Add(24 * time.Hour), CreatedAt: time.Now()},
			},
			TotalItems: 3,
		}, nil).Once()

		resp, err := client.ListInvites(context.Background(), &authv1.ListInvitesRequest{
			Page:         1,
			PageSize:     10,
			CreatedAfter: &after,
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(resp.GetInvites()), 3)
	})
}

func TestAuthServer_CancelInvite(t *testing.T) {
	t.Run("cancel pending invite succeeds", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)

		// InviteUser expectations
		userRepo.On("GetByEmail", mock.Anything, "cancelme@example.com").Return(nil, nil).Once()
		inviteRepo.On("ExistsPendingByEmail", mock.Anything, "cancelme@example.com").Return(false, nil).Once()
		inviteRepo.On("Create", mock.Anything, mock.MatchedBy(func(inv *model.InviteToken) bool {
			return inv.Email == "cancelme@example.com"
		})).Return(nil).Once()

		inviteResp, err := client.InviteUser(context.Background(), &authv1.InviteUserRequest{
			Email:     "cancelme@example.com",
			InvitedBy: "admin-123",
		})
		require.NoError(t, err)
		require.NotEmpty(t, inviteResp.GetToken())

		// CancelInvite expectations
		inviteRepo.On("GetByToken", mock.Anything, inviteResp.GetToken()).Return(&model.InviteToken{
			ID:        "inv-1",
			Email:     "cancelme@example.com",
			Token:     inviteResp.GetToken(),
			InvitedBy: "admin-123",
			Status:    model.StatusPending,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now(),
		}, nil).Once()
		inviteRepo.On("Cancel", mock.Anything, inviteResp.GetToken()).Return(nil).Once()

		cancelResp, err := client.CancelInvite(context.Background(), &authv1.CancelInviteRequest{
			Token: inviteResp.GetToken(),
		})
		require.NoError(t, err)
		assert.True(t, cancelResp.GetSuccess())
	})

	t.Run("cancel non-existent invite returns error", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)

		inviteRepo.On("GetByToken", mock.Anything, "nonexistent-token").Return(nil, nil).Once()

		_, err := client.CancelInvite(context.Background(), &authv1.CancelInviteRequest{
			Token: "nonexistent-token",
		})
		assert.Error(t, err)
	})

	t.Run("cancel empty token returns error", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)

		_, err := client.CancelInvite(context.Background(), &authv1.CancelInviteRequest{
			Token: "",
		})
		assert.Error(t, err)
	})
}

func TestAuthServer_GetPublicKey(t *testing.T) {
	t.Run("returns PEM public key", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewAuthServiceClient(conn)

		resp, err := client.GetPublicKey(context.Background(), &authv1.GetPublicKeyRequest{})
		require.NoError(t, err)
		assert.Contains(t, resp.GetPublicKeyPem(), "-----BEGIN PUBLIC KEY-----")
	})
}
