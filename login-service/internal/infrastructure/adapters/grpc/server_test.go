package grpc

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"

	"github.com/beabys/wms/login-service/internal/application/auth/usecase"
	"github.com/beabys/wms/login-service/internal/domain/auth/model"
)

var errNotFound = errors.New("not found")

// bufConnRepo is a mock UserRepository for in-process gRPC testing.
type bufConnRepo struct {
	users         map[string]*model.User
	refreshTokens map[string]string
}

func newBufConnRepo() *bufConnRepo {
	return &bufConnRepo{
		users:         make(map[string]*model.User),
		refreshTokens: make(map[string]string),
	}
}

func (m *bufConnRepo) GetByID(_ context.Context, id string) (*model.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errNotFound
	}
	return u, nil
}

func (m *bufConnRepo) GetByEmail(_ context.Context, email model.Email) (*model.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errNotFound
}

func (m *bufConnRepo) Save(_ context.Context, user *model.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *bufConnRepo) List(_ context.Context, limit, offset int32) ([]*model.User, error) {
	var result []*model.User
	for _, u := range m.users {
		result = append(result, u)
	}
	return result, nil
}

func (m *bufConnRepo) Update(_ context.Context, user *model.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *bufConnRepo) CreateRefreshToken(_ context.Context, userID, token string) error {
	m.refreshTokens[token] = userID
	return nil
}

func (m *bufConnRepo) ValidateRefreshToken(_ context.Context, token string) (string, error) {
	uid, ok := m.refreshTokens[token]
	if !ok {
		return "", errNotFound
	}
	return uid, nil
}

func (m *bufConnRepo) DeleteRefreshToken(_ context.Context, token string) error {
	delete(m.refreshTokens, token)
	return nil
}

func mustGenerateKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return key, &key.PublicKey
}

func setupTestGRPCServer(t *testing.T) (authv1.AuthServiceClient, authv1.UserServiceClient, *bufConnRepo) {
	t.Helper()

	repo := newBufConnRepo()
	privKey, pubKey := mustGenerateKeys(t)
	logger := zap.NewNop()

	// Create usecase
	authService := usecase.NewAuthService(repo, privKey, pubKey, 15*time.Minute, 7*24*time.Hour)

	// Create gRPC servers
	authServer := NewAuthServer(authService, logger, "")
	userServer := NewUserServer(authService, logger)

	// In-process gRPC setup
	bufSize := 1024 * 1024
	listener := bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	authv1.RegisterAuthServiceServer(srv, authServer)
	authv1.RegisterUserServiceServer(srv, userServer)

	go func() {
		require.NoError(t, srv.Serve(listener))
	}()
	t.Cleanup(func() { srv.Stop() })

	// Client
	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	authClient := authv1.NewAuthServiceClient(conn)
	userClient := authv1.NewUserServiceClient(conn)

	return authClient, userClient, repo
}

func TestAuthServer_Login_Success(t *testing.T) {
	authClient, _, repo := setupTestGRPCServer(t)

	// Seed user
	email, _ := model.NewEmail("grpc@example.com")
	passHash, _ := model.NewPasswordHash("testPass123")
	repo.Save(context.Background(), &model.User{
		ID:           "grpc-user-1",
		Email:        email,
		PasswordHash: passHash,
		Role:         model.RoleCustomer,
		Status:       model.UserStatusActive,
	})

	resp, err := authClient.Login(context.Background(), &authv1.LoginRequest{
		Email:    "grpc@example.com",
		Password: "testPass123",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.GetAccessToken())
	require.NotEmpty(t, resp.GetRefreshToken())
	require.NotNil(t, resp.GetUser())
	require.Equal(t, "grpc-user-1", resp.GetUser().GetId())
}

func TestAuthServer_Login_InvalidCredentials(t *testing.T) {
	authClient, _, _ := setupTestGRPCServer(t)

	_, err := authClient.Login(context.Background(), &authv1.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "wrongPass123",
	})

	require.Error(t, err)
}

func TestAuthServer_RefreshToken_Success(t *testing.T) {
	authClient, _, repo := setupTestGRPCServer(t)

	// Seed user and login to get refresh token
	email, _ := model.NewEmail("refresh-grpc@example.com")
	passHash, _ := model.NewPasswordHash("testPass123")
	repo.Save(context.Background(), &model.User{
		ID:           "refresh-grpc-1",
		Email:        email,
		PasswordHash: passHash,
		Role:         model.RoleCustomer,
		Status:       model.UserStatusActive,
	})

	loginResp, err := authClient.Login(context.Background(), &authv1.LoginRequest{
		Email:    "refresh-grpc@example.com",
		Password: "testPass123",
	})
	require.NoError(t, err)
	require.NotEmpty(t, loginResp.GetRefreshToken())

	// Now refresh
	refreshResp, err := authClient.RefreshToken(context.Background(), &authv1.RefreshTokenRequest{
		RefreshToken: loginResp.GetRefreshToken(),
	})

	require.NoError(t, err)
	require.NotEmpty(t, refreshResp.GetAccessToken())
	require.NotEmpty(t, refreshResp.GetRefreshToken())
}

func TestUserServer_CreateUser_Success(t *testing.T) {
	_, userClient, _ := setupTestGRPCServer(t)

	resp, err := userClient.CreateUser(context.Background(), &authv1.CreateUserRequest{
		Email:    "newuser@example.com",
		Password: "securePass123",
		Role:     "customer",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetUser())
	require.Equal(t, "newuser@example.com", resp.GetUser().GetEmail())
	require.Equal(t, "customer", resp.GetUser().GetRole())
}

func TestAuthServer_GetPublicKey(t *testing.T) {
	authClient, _, _ := setupTestGRPCServer(t)

	resp, err := authClient.GetPublicKey(context.Background(), &authv1.GetPublicKeyRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "rsa-1", resp.GetKeyId())
}

func TestUserServer_ListUsers(t *testing.T) {
	_, userClient, repo := setupTestGRPCServer(t)

	// Seed users
	for i := 0; i < 3; i++ {
		email, _ := model.NewEmail("user" + string(rune('0'+i)) + "@example.com")
		passHash, _ := model.NewPasswordHash("password123")
		repo.Save(context.Background(), &model.User{
			ID:           "user-" + string(rune('0'+i)),
			Email:        email,
			PasswordHash: passHash,
			Role:         model.RoleCustomer,
			Status:       model.UserStatusActive,
		})
	}

	resp, err := userClient.ListUsers(context.Background(), &authv1.ListUsersRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.GetUsers(), 3)
}
