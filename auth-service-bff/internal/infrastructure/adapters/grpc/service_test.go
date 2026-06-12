package grpcdapter

import (
	"context"
	"net"
	"testing"
	"time"

	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// testAuthServer implements authv1.AuthServiceServer for testing.
type testAuthServer struct {
	authv1.UnimplementedAuthServiceServer
	loginFn        func(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error)
	refreshTokenFn func(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error)
	validateTokenFn func(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error)
	getPublicKeyFn func(ctx context.Context, req *authv1.GetPublicKeyRequest) (*authv1.GetPublicKeyResponse, error)
	listInvitesFn  func(ctx context.Context, req *authv1.ListInvitesRequest) (*authv1.ListInvitesResponse, error)
	cancelInviteFn func(ctx context.Context, req *authv1.CancelInviteRequest) (*authv1.CancelInviteResponse, error)
}

func (s *testAuthServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if s.loginFn != nil {
		return s.loginFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *testAuthServer) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	if s.refreshTokenFn != nil {
		return s.refreshTokenFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *testAuthServer) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	if s.validateTokenFn != nil {
		return s.validateTokenFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *testAuthServer) GetPublicKey(ctx context.Context, req *authv1.GetPublicKeyRequest) (*authv1.GetPublicKeyResponse, error) {
	if s.getPublicKeyFn != nil {
		return s.getPublicKeyFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *testAuthServer) ListInvites(ctx context.Context, req *authv1.ListInvitesRequest) (*authv1.ListInvitesResponse, error) {
	if s.listInvitesFn != nil {
		return s.listInvitesFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *testAuthServer) CancelInvite(ctx context.Context, req *authv1.CancelInviteRequest) (*authv1.CancelInviteResponse, error) {
	if s.cancelInviteFn != nil {
		return s.cancelInviteFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// testUserServer implements authv1.UserServiceServer for testing.
type testUserServer struct {
	authv1.UnimplementedUserServiceServer
	createUserFn func(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error)
	getUserFn    func(ctx context.Context, req *authv1.GetUserRequest) (*authv1.GetUserResponse, error)
	updateUserFn func(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.UpdateUserResponse, error)
	listUsersFn  func(ctx context.Context, req *authv1.ListUsersRequest) (*authv1.ListUsersResponse, error)
	deleteUserFn func(ctx context.Context, req *authv1.DeleteUserRequest) (*authv1.DeleteUserResponse, error)
	assignRoleFn func(ctx context.Context, req *authv1.AssignRoleRequest) (*authv1.AssignRoleResponse, error)
	listRolesFn  func(ctx context.Context, req *authv1.ListRolesRequest) (*authv1.ListRolesResponse, error)
	createRoleFn func(ctx context.Context, req *authv1.CreateRoleRequest) (*authv1.CreateRoleResponse, error)
}

func (s *testUserServer) CreateUser(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
	if s.createUserFn != nil {
		return s.createUserFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (s *testUserServer) GetUser(ctx context.Context, req *authv1.GetUserRequest) (*authv1.GetUserResponse, error) {
	if s.getUserFn != nil {
		return s.getUserFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (s *testUserServer) UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.UpdateUserResponse, error) {
	if s.updateUserFn != nil {
		return s.updateUserFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (s *testUserServer) ListUsers(ctx context.Context, req *authv1.ListUsersRequest) (*authv1.ListUsersResponse, error) {
	if s.listUsersFn != nil {
		return s.listUsersFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (s *testUserServer) DeleteUser(ctx context.Context, req *authv1.DeleteUserRequest) (*authv1.DeleteUserResponse, error) {
	if s.deleteUserFn != nil {
		return s.deleteUserFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (s *testUserServer) AssignRole(ctx context.Context, req *authv1.AssignRoleRequest) (*authv1.AssignRoleResponse, error) {
	if s.assignRoleFn != nil {
		return s.assignRoleFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (s *testUserServer) ListRoles(ctx context.Context, req *authv1.ListRolesRequest) (*authv1.ListRolesResponse, error) {
	if s.listRolesFn != nil {
		return s.listRolesFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (s *testUserServer) CreateRole(ctx context.Context, req *authv1.CreateRoleRequest) (*authv1.CreateRoleResponse, error) {
	if s.createRoleFn != nil {
		return s.createRoleFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// setupTestServer starts a local gRPC server and returns a client connected to it.
func setupTestServer(t *testing.T, authSrv *testAuthServer, userSrv *testUserServer) (*AuthClient, func()) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	srv := grpc.NewServer()
	if authSrv != nil {
		authv1.RegisterAuthServiceServer(srv, authSrv)
	}
	if userSrv != nil {
		authv1.RegisterUserServiceServer(srv, userSrv)
	}

	go func() {
		srv.Serve(listener) //nolint:errcheck
	}()

	client, err := NewAuthClient(listener.Addr().String())
	require.NoError(t, err)

	cleanup := func() {
		client.Close()
		srv.GracefulStop()
	}

	return client, cleanup
}

func TestAuthServiceAdapterLogin(t *testing.T) {
	authSrv := &testAuthServer{
		loginFn: func(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
			return &authv1.LoginResponse{
				AccessToken:  "access-token-123",
				RefreshToken: "refresh-token-456",
				ExpiresIn:    3600,
			}, nil
		},
	}
	client, cleanup := setupTestServer(t, authSrv, nil)
	defer cleanup()

	adapter := NewAuthServiceAdapter(client)
	resp, err := adapter.Login(context.Background(), "user@example.com", "password")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "access-token-123", resp.AccessToken)
	assert.Equal(t, "refresh-token-456", resp.RefreshToken)
	assert.Equal(t, int64(3600), resp.ExpiresIn)
}

func TestAuthServiceAdapterLoginError(t *testing.T) {
	authSrv := &testAuthServer{
		loginFn: func(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
			return nil, status.Error(codes.Unauthenticated, "bad credentials")
		},
	}
	client, cleanup := setupTestServer(t, authSrv, nil)
	defer cleanup()

	adapter := NewAuthServiceAdapter(client)
	resp, err := adapter.Login(context.Background(), "bad@user.com", "wrong")
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, ErrUnauthenticated)
}

func TestAuthServiceAdapterRefreshToken(t *testing.T) {
	authSrv := &testAuthServer{
		refreshTokenFn: func(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
			return &authv1.RefreshTokenResponse{
				AccessToken:  "new-access",
				RefreshToken: "new-refresh",
				ExpiresIn:    1800,
			}, nil
		},
	}
	client, cleanup := setupTestServer(t, authSrv, nil)
	defer cleanup()

	adapter := NewAuthServiceAdapter(client)
	resp, err := adapter.RefreshToken(context.Background(), "valid-refresh-token")
	require.NoError(t, err)
	assert.Equal(t, "new-access", resp.AccessToken)
}

func TestAuthServiceAdapterRefreshTokenError(t *testing.T) {
	authSrv := &testAuthServer{
		refreshTokenFn: func(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
			return nil, status.Error(codes.Unauthenticated, "invalid refresh")
		},
	}
	client, cleanup := setupTestServer(t, authSrv, nil)
	defer cleanup()

	adapter := NewAuthServiceAdapter(client)
	resp, err := adapter.RefreshToken(context.Background(), "bad-token")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestAuthServiceAdapterValidateToken(t *testing.T) {
	authSrv := &testAuthServer{
		validateTokenFn: func(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
			return &authv1.ValidateTokenResponse{
				Valid:  true,
				UserId: "user-1",
				Role:   "admin",
			}, nil
		},
	}
	client, cleanup := setupTestServer(t, authSrv, nil)
	defer cleanup()

	adapter := NewAuthServiceAdapter(client)
	resp, err := adapter.ValidateToken(context.Background(), "valid-jwt")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.GetValid())
	assert.Equal(t, "user-1", resp.GetUserId())
	assert.Equal(t, "admin", resp.GetRole())
}

func TestAuthServiceAdapterValidateTokenError(t *testing.T) {
	authSrv := &testAuthServer{
		validateTokenFn: func(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		},
	}
	client, cleanup := setupTestServer(t, authSrv, nil)
	defer cleanup()

	adapter := NewAuthServiceAdapter(client)
	resp, err := adapter.ValidateToken(context.Background(), "bad-jwt")
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, ErrUnauthenticated)
}

func TestAuthServiceAdapterCreateUser(t *testing.T) {
	userSrv := &testUserServer{
		createUserFn: func(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
			return &authv1.CreateUserResponse{
				User: &authv1.User{
					Id:     "new-user-1",
					Email:  req.Email,
					Name:   req.Name,
					Role:   req.Role,
					Active: true,
					CreatedAt: &commonv1.Timestamp{Seconds: time.Now().Unix()},
				},
			}, nil
		},
	}
	client, cleanup := setupTestServer(t, nil, userSrv)
	defer cleanup()

	adapter := NewAuthServiceAdapter(client)
	resp, err := adapter.CreateUser(context.Background(), &authv1.CreateUserRequest{
		Email:    "new@example.com",
		Password: "secure-pass",
		Name:     "New User",
		Role:     "viewer",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "new-user-1", resp.GetUser().GetId())
	assert.Equal(t, "new@example.com", resp.GetUser().GetEmail())
}

func TestAuthServiceAdapterListUsers(t *testing.T) {
	userSrv := &testUserServer{
		listUsersFn: func(ctx context.Context, req *authv1.ListUsersRequest) (*authv1.ListUsersResponse, error) {
			return &authv1.ListUsersResponse{
				Users: []*authv1.User{
					{Id: "u1", Email: "a@b.com", Name: "A", Role: "admin", Active: true},
					{Id: "u2", Email: "c@d.com", Name: "C", Role: "viewer", Active: true},
				},
				Pagination: &commonv1.Pagination{
					Page: 1, PageSize: 10, Total: 2,
				},
			}, nil
		},
	}
	client, cleanup := setupTestServer(t, nil, userSrv)
	defer cleanup()

	adapter := NewAuthServiceAdapter(client)
	resp, err := adapter.ListUsers(context.Background(), &authv1.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Len(t, resp.GetUsers(), 2)
	assert.Equal(t, int32(1), resp.GetPagination().GetPage())
	assert.Equal(t, int32(2), resp.GetPagination().GetTotal())
}

func TestAuthServiceAdapterGetUser(t *testing.T) {
	userSrv := &testUserServer{
		getUserFn: func(ctx context.Context, req *authv1.GetUserRequest) (*authv1.GetUserResponse, error) {
			return &authv1.GetUserResponse{
				User: &authv1.User{
					Id:     req.Id,
					Email:  "found@example.com",
					Name:   "Found User",
					Role:   "admin",
					Active: true,
				},
			}, nil
		},
	}
	client, cleanup := setupTestServer(t, nil, userSrv)
	defer cleanup()

	adapter := NewAuthServiceAdapter(client)
	resp, err := adapter.GetUser(context.Background(), "user-42")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "user-42", resp.GetUser().GetId())
	assert.Equal(t, "found@example.com", resp.GetUser().GetEmail())
}

func TestAuthServiceAdapterGetUserNotFound(t *testing.T) {
	userSrv := &testUserServer{
		getUserFn: func(ctx context.Context, req *authv1.GetUserRequest) (*authv1.GetUserResponse, error) {
			return nil, status.Error(codes.NotFound, "user not found")
		},
	}
	client, cleanup := setupTestServer(t, nil, userSrv)
	defer cleanup()

	adapter := NewAuthServiceAdapter(client)
	resp, err := adapter.GetUser(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestAuthServiceAdapterUpdateUser(t *testing.T) {
	userSrv := &testUserServer{
		updateUserFn: func(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.UpdateUserResponse, error) {
			return &authv1.UpdateUserResponse{
				User: &authv1.User{
					Id:     req.Id,
					Email:  "updated@example.com",
					Name:   req.Name,
					Role:   req.Role,
					Active: req.Active.GetValue(),
				},
			}, nil
		},
	}
	client, cleanup := setupTestServer(t, nil, userSrv)
	defer cleanup()

	adapter := NewAuthServiceAdapter(client)
	resp, err := adapter.UpdateUser(context.Background(), &authv1.UpdateUserRequest{
		Id:   "user-1",
		Name: "Updated Name",
		Role: "manager",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "Updated Name", resp.GetUser().GetName())
	assert.Equal(t, "manager", resp.GetUser().GetRole())
}

func TestAuthServiceAdapterDeleteUser(t *testing.T) {
	userSrv := &testUserServer{
		deleteUserFn: func(ctx context.Context, req *authv1.DeleteUserRequest) (*authv1.DeleteUserResponse, error) {
			return &authv1.DeleteUserResponse{}, nil
		},
	}
	client, cleanup := setupTestServer(t, nil, userSrv)
	defer cleanup()

	adapter := NewAuthServiceAdapter(client)
	_, err := adapter.DeleteUser(context.Background(), "user-1")
	assert.NoError(t, err)
}

func TestAuthServiceAdapterDeleteUserError(t *testing.T) {
	userSrv := &testUserServer{
		deleteUserFn: func(ctx context.Context, req *authv1.DeleteUserRequest) (*authv1.DeleteUserResponse, error) {
			return nil, status.Error(codes.PermissionDenied, "no permission")
		},
	}
	client, cleanup := setupTestServer(t, nil, userSrv)
	defer cleanup()

	adapter := NewAuthServiceAdapter(client)
	_, err := adapter.DeleteUser(context.Background(), "user-1")
	assert.ErrorIs(t, err, ErrPermissionDenied)
}

func TestAuthServiceAdapterAssignRole(t *testing.T) {
	userSrv := &testUserServer{
		assignRoleFn: func(ctx context.Context, req *authv1.AssignRoleRequest) (*authv1.AssignRoleResponse, error) {
			return &authv1.AssignRoleResponse{}, nil
		},
	}
	client, cleanup := setupTestServer(t, nil, userSrv)
	defer cleanup()

	adapter := NewAuthServiceAdapter(client)
	_, err := adapter.AssignRole(context.Background(), "user-1", "admin")
	assert.NoError(t, err)
}

func TestAuthServiceAdapterListRoles(t *testing.T) {
	userSrv := &testUserServer{
		listRolesFn: func(ctx context.Context, req *authv1.ListRolesRequest) (*authv1.ListRolesResponse, error) {
			return &authv1.ListRolesResponse{
				Roles: []*authv1.Role{
					{Id: "r1", Name: "admin", Description: "Admin role"},
				},
			}, nil
		},
	}
	client, cleanup := setupTestServer(t, nil, userSrv)
	defer cleanup()

	adapter := NewAuthServiceAdapter(client)
	resp, err := adapter.ListRoles(context.Background(), &authv1.ListRolesRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Len(t, resp.GetRoles(), 1)
	assert.Equal(t, "admin", resp.GetRoles()[0].GetName())
}

func TestAuthServiceAdapterCreateRole(t *testing.T) {
	userSrv := &testUserServer{
		createRoleFn: func(ctx context.Context, req *authv1.CreateRoleRequest) (*authv1.CreateRoleResponse, error) {
			return &authv1.CreateRoleResponse{
				Role: &authv1.Role{
					Id:          "new-role-1",
					Name:        req.Name,
					Description: req.Description,
				},
			}, nil
		},
	}
	client, cleanup := setupTestServer(t, nil, userSrv)
	defer cleanup()

	adapter := NewAuthServiceAdapter(client)
	resp, err := adapter.CreateRole(context.Background(), &authv1.CreateRoleRequest{
		Name: "editor", Description: "Can edit",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "new-role-1", resp.GetRole().GetId())
	assert.Equal(t, "editor", resp.GetRole().GetName())
}

func TestAuthServiceAdapterListInvites(t *testing.T) {
	t.Run("success with no filters", func(t *testing.T) {
		authSrv := &testAuthServer{
			listInvitesFn: func(ctx context.Context, req *authv1.ListInvitesRequest) (*authv1.ListInvitesResponse, error) {
				assert.Equal(t, int32(1), req.GetPage())
				assert.Equal(t, int32(20), req.GetPageSize())
				return &authv1.ListInvitesResponse{
					Invites: []*authv1.InviteEntry{
						{Id: "1", Email: "test@example.com", Token: "abc", InvitedBy: "admin", Status: "pending", ExpiresAt: 1900000000, CreatedAt: 1800000000},
					},
					Pagination: &commonv1.Pagination{Page: 1, PageSize: 20, Total: 1, TotalPages: 1},
				}, nil
			},
		}
		client, cleanup := setupTestServer(t, authSrv, nil)
		defer cleanup()

		adapter := NewAuthServiceAdapter(client)
		resp, err := adapter.ListInvites(context.Background(), &authv1.ListInvitesRequest{
			Page:     1,
			PageSize: 20,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Len(t, resp.GetInvites(), 1)
		assert.Equal(t, "test@example.com", resp.GetInvites()[0].GetEmail())
		assert.Equal(t, int32(1), resp.GetPagination().GetPage())
	})

	t.Run("with filters", func(t *testing.T) {
		authSrv := &testAuthServer{
			listInvitesFn: func(ctx context.Context, req *authv1.ListInvitesRequest) (*authv1.ListInvitesResponse, error) {
				assert.Equal(t, "used", req.GetStatus())
				assert.False(t, req.GetExpired())
				assert.Equal(t, int64(1000000), req.GetCreatedAfter())
				return &authv1.ListInvitesResponse{
					Invites:    []*authv1.InviteEntry{},
					Pagination: &commonv1.Pagination{Page: 1, PageSize: 20, Total: 0, TotalPages: 0},
				}, nil
			},
		}
		client, cleanup := setupTestServer(t, authSrv, nil)
		defer cleanup()

		adapter := NewAuthServiceAdapter(client)
		status := "used"
		expired := false
		after := int64(1000000)
		resp, err := adapter.ListInvites(context.Background(), &authv1.ListInvitesRequest{
			Status:       &status,
			Expired:      &expired,
			CreatedAfter: &after,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Empty(t, resp.GetInvites())
	})

	t.Run("error mapping", func(t *testing.T) {
		authSrv := &testAuthServer{
			listInvitesFn: func(ctx context.Context, req *authv1.ListInvitesRequest) (*authv1.ListInvitesResponse, error) {
				return nil, status.Error(codes.Internal, "internal error")
			},
		}
		client, cleanup := setupTestServer(t, authSrv, nil)
		defer cleanup()

		adapter := NewAuthServiceAdapter(client)
		resp, err := adapter.ListInvites(context.Background(), &authv1.ListInvitesRequest{})
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestAuthServiceAdapterCancelInvite(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		authSrv := &testAuthServer{
			cancelInviteFn: func(ctx context.Context, req *authv1.CancelInviteRequest) (*authv1.CancelInviteResponse, error) {
				assert.Equal(t, "cancel-token", req.GetToken())
				return &authv1.CancelInviteResponse{Success: true}, nil
			},
		}
		client, cleanup := setupTestServer(t, authSrv, nil)
		defer cleanup()

		adapter := NewAuthServiceAdapter(client)
		resp, err := adapter.CancelInvite(context.Background(), "cancel-token")
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.True(t, resp.Success)
	})

	t.Run("error mapping", func(t *testing.T) {
		authSrv := &testAuthServer{
			cancelInviteFn: func(ctx context.Context, req *authv1.CancelInviteRequest) (*authv1.CancelInviteResponse, error) {
				return nil, status.Error(codes.InvalidArgument, "cannot cancel")
			},
		}
		client, cleanup := setupTestServer(t, authSrv, nil)
		defer cleanup()

		adapter := NewAuthServiceAdapter(client)
		resp, err := adapter.CancelInvite(context.Background(), "bad-token")
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestGRPCDialErrorPassthrough(t *testing.T) {
	// Test that NewAuthClient with unreachable target still creates client (lazy dial)
	client, err := NewAuthClient("127.0.0.1:1")
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	// Actual RPC call will fail — test that adapter wraps the error
	adapter := NewAuthServiceAdapter(client)
	resp, err := adapter.Login(context.Background(), "test@test.com", "pass")
	assert.Error(t, err)
	assert.Nil(t, resp)
}
