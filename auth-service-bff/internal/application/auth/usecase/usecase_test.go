package usecase

import (
	"context"
	"testing"

	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ucmocks "github.com/beabys/wms/auth-service-bff/mocks/application/auth/usecase"
	"github.com/beabys/wms/auth-service-bff/internal/domain/model"
	"github.com/beabys/wms/auth-service-bff/internal/infrastructure/adapters/http/context"
	"github.com/beabys/wms/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func setupUseCase(t *testing.T, mockSvc *ucmocks.GrpcClient) *AuthUseCase {
	t.Helper()
	log, err := logger.NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	require.NoError(t, err)
	return New(log, mockSvc)
}

// --- Login tests ---

func TestLogin_Success(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	mockSvc.EXPECT().Login(mock.Anything, "user@example.com", "password123").
		Return(&authv1.LoginResponse{AccessToken: "mock-access", RefreshToken: "mock-refresh", ExpiresIn: 3600}, nil)
	uc := setupUseCase(t, mockSvc)
	resp, err := uc.Login(context.Background(), "user@example.com", "password123")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "mock-access", resp.AccessToken)
	assert.Equal(t, int64(3600), resp.ExpiresIn)
}

func TestLogin_InvalidEmail(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	uc := setupUseCase(t, mockSvc)
	_, err := uc.Login(context.Background(), "invalid", "password123")
	assert.Error(t, err)
}

func TestLogin_EmptyPassword(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	uc := setupUseCase(t, mockSvc)
	_, err := uc.Login(context.Background(), "user@example.com", "")
	assert.Error(t, err)
}

// --- GetMe tests ---

func TestGetMe_Success(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	mockSvc.EXPECT().ValidateToken(mock.Anything, "valid-jwt").
		Return(&authv1.ValidateTokenResponse{Valid: true, UserId: "user-1", Role: "admin"}, nil)
	mockSvc.EXPECT().GetUser(mock.Anything, "user-1").
		Return(&authv1.GetUserResponse{User: &authv1.User{Id: "user-1", Email: "found@example.com", Name: "Found", Role: "viewer", Active: true}}, nil)
	uc := setupUseCase(t, mockSvc)
	ctx := context.WithValue(context.Background(), httpctx.ContextKeyJWT, "valid-jwt")
	resp, err := uc.GetMe(ctx)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "user-1", resp.ID)
	assert.Equal(t, "found@example.com", resp.Email)
}

func TestGetMe_NoToken(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	uc := setupUseCase(t, mockSvc)
	_, err := uc.GetMe(context.Background())
	assert.Error(t, err)
}

// --- CreateUser tests ---

func TestCreateUser_Success(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	cid := "cust-1"
	mockSvc.EXPECT().CreateUser(mock.Anything, mock.MatchedBy(func(req *authv1.CreateUserRequest) bool {
		return req.Email == "new@example.com" && req.Name == "New User" && req.Role == "viewer"
	})).
		Return(&authv1.CreateUserResponse{
			User: &authv1.User{Id: "new-user", Email: "new@example.com", Name: "New User", Role: "viewer", Active: true},
		}, nil)
	uc := setupUseCase(t, mockSvc)
	resp, err := uc.CreateUser(context.Background(), &model.CreateUserRequest{
		Email:      "new@example.com",
		Password:   "password123",
		Name:       "New User",
		Role:       "viewer",
		CustomerID: &cid,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "new-user", resp.ID)
	assert.Equal(t, "new@example.com", resp.Email)
}

func TestCreateUser_ValidationError(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	uc := setupUseCase(t, mockSvc)
	_, err := uc.CreateUser(context.Background(), &model.CreateUserRequest{
		Email:    "bad",
		Password: "short",
		Name:     "",
		Role:     "",
	})
	assert.Error(t, err)
}

// --- InviteUser tests ---

func TestInviteUser_Success(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	mockSvc.EXPECT().InviteUser(mock.Anything, "invite@example.com", "admin").
		Return(&authv1.InviteUserResponse{Token: "invite-token", InviteLink: "https://example.com/register?token=invite-token", ExpiresAt: 1800000000}, nil)
	uc := setupUseCase(t, mockSvc)
	resp, err := uc.InviteUser(context.Background(), "invite@example.com", "admin")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "invite-token", resp.Token)
}

func TestInviteUser_InvalidEmail(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	uc := setupUseCase(t, mockSvc)
	_, err := uc.InviteUser(context.Background(), "", "admin")
	assert.Error(t, err)
}

// --- ListInvites tests ---

func TestListInvites_Success(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	mockSvc.EXPECT().ListInvites(mock.Anything, mock.AnythingOfType("*authv1.ListInvitesRequest")).
		Return(&authv1.ListInvitesResponse{
			Invites:    []*authv1.InviteEntry{{Id: "1", Email: "invited@example.com", Token: "abc", InvitedBy: "admin", Status: "pending"}},
			Pagination: &commonv1.Pagination{Page: 1, PageSize: 20, Total: 1},
		}, nil)
	uc := setupUseCase(t, mockSvc)
	page := 1
	pageSize := 20
	resp, err := uc.ListInvites(context.Background(), nil, nil, nil, nil, &page, &pageSize)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Len(t, resp.Invites, 1)
	assert.Equal(t, 1, resp.Pagination.Page)
}

// --- CancelInvite tests ---

func TestCancelInvite_Success(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	mockSvc.EXPECT().CancelInvite(mock.Anything, "token-123").
		Return(&authv1.CancelInviteResponse{Success: true}, nil)
	uc := setupUseCase(t, mockSvc)
	resp, err := uc.CancelInvite(context.Background(), "token-123")
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestCancelInvite_EmptyToken(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	uc := setupUseCase(t, mockSvc)
	_, err := uc.CancelInvite(context.Background(), "")
	assert.Error(t, err)
}

// --- ListUsers tests ---

func TestListUsers_Success(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	mockSvc.EXPECT().ListUsers(mock.Anything, mock.AnythingOfType("*authv1.ListUsersRequest")).
		Return(&authv1.ListUsersResponse{
			Users:      []*authv1.User{{Id: "u1", Email: "a@b.com", Name: "A", Role: "viewer", Active: true}},
			Pagination: &commonv1.Pagination{Page: 1, PageSize: 10, Total: 1},
		}, nil)
	uc := setupUseCase(t, mockSvc)
	page := 1
	pageSize := 10
	resp, err := uc.ListUsers(context.Background(), nil, nil, &page, &pageSize)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Len(t, resp.Users, 1)
	assert.Equal(t, 1, resp.Pagination.Page)
}

// --- GetUser tests ---

func TestGetUser_Success(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	mockSvc.EXPECT().GetUser(mock.Anything, "user-42").
		Return(&authv1.GetUserResponse{User: &authv1.User{Id: "user-42", Email: "found@example.com", Name: "Found", Role: "viewer", Active: true}}, nil)
	uc := setupUseCase(t, mockSvc)
	resp, err := uc.GetUser(context.Background(), "user-42")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "user-42", resp.ID)
	assert.Equal(t, "found@example.com", resp.Email)
}

// --- UpdateUser tests ---

func TestUpdateUser_Success(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	name := "Updated"
	role := "manager"
	mockSvc.EXPECT().UpdateUser(mock.Anything, mock.MatchedBy(func(req *authv1.UpdateUserRequest) bool {
		return req.Id == "user-1"
	})).
		Return(&authv1.UpdateUserResponse{User: &authv1.User{Id: "user-1", Email: "updated@example.com", Name: name, Role: role, Active: true}}, nil)
	uc := setupUseCase(t, mockSvc)
	active := true
	resp, err := uc.UpdateUser(context.Background(), "user-1", &model.UpdateUserRequest{
		Name: &name, Role: &role, Active: &active,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "Updated", resp.Name)
}

// --- DeleteUser tests ---

func TestDeleteUser_Success(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	mockSvc.EXPECT().DeleteUser(mock.Anything, "user-1").Return(&authv1.DeleteUserResponse{}, nil)
	uc := setupUseCase(t, mockSvc)
	err := uc.DeleteUser(context.Background(), "user-1")
	assert.NoError(t, err)
}

// --- AssignRole tests ---

func TestAssignRole_Success(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	mockSvc.EXPECT().AssignRole(mock.Anything, "user-1", "admin").Return(&authv1.AssignRoleResponse{}, nil)
	uc := setupUseCase(t, mockSvc)
	err := uc.AssignRole(context.Background(), "user-1", "admin")
	assert.NoError(t, err)
}

func TestAssignRole_EmptyRole(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	uc := setupUseCase(t, mockSvc)
	err := uc.AssignRole(context.Background(), "user-1", "")
	assert.Error(t, err)
}

// --- ListRoles tests ---

func TestListRoles_Success(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	mockSvc.EXPECT().ListRoles(mock.Anything, mock.AnythingOfType("*authv1.ListRolesRequest")).
		Return(&authv1.ListRolesResponse{Roles: []*authv1.Role{{Id: "r1", Name: "admin", Description: "Admin role"}}}, nil)
	uc := setupUseCase(t, mockSvc)
	resp, err := uc.ListRoles(context.Background())
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Len(t, resp.Roles, 1)
	assert.Equal(t, "admin", resp.Roles[0].Name)
}

// --- CreateRole tests ---

func TestCreateRole_Success(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	mockSvc.EXPECT().CreateRole(mock.Anything, mock.AnythingOfType("*authv1.CreateRoleRequest")).
		Return(&authv1.CreateRoleResponse{Role: &authv1.Role{Id: "new-role", Name: "editor", Description: "Can edit"}}, nil)
	uc := setupUseCase(t, mockSvc)
	resp, err := uc.CreateRole(context.Background(), &model.CreateRoleRequest{
		Name: "editor", Description: "Can edit",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "new-role", resp.ID)
	assert.Equal(t, "editor", resp.Name)
}

func TestCreateRole_EmptyName(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	uc := setupUseCase(t, mockSvc)
	_, err := uc.CreateRole(context.Background(), &model.CreateRoleRequest{
		Name: "", Description: "no name",
	})
	assert.Error(t, err)
}

// --- Logout tests ---

func TestLogout_Success(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	mockSvc.EXPECT().Logout(mock.Anything, "refresh-token").Return(&authv1.LogoutResponse{}, nil)
	uc := setupUseCase(t, mockSvc)
	err := uc.Logout(context.Background(), "refresh-token")
	assert.NoError(t, err)
}

// --- RefreshToken tests ---

func TestRefreshToken_Success(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	mockSvc.EXPECT().RefreshToken(mock.Anything, "valid-refresh").
		Return(&authv1.RefreshTokenResponse{AccessToken: "new-access", RefreshToken: "new-refresh", ExpiresIn: 1800}, nil)
	uc := setupUseCase(t, mockSvc)
	resp, err := uc.RefreshToken(context.Background(), "valid-refresh")
	require.NoError(t, err)
	assert.Equal(t, "new-access", resp.AccessToken)
}

func TestRefreshToken_Empty(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	uc := setupUseCase(t, mockSvc)
	_, err := uc.RefreshToken(context.Background(), "")
	assert.Error(t, err)
}

// --- gRPC error passthrough ---

func TestLogin_GrpcError(t *testing.T) {
	mockSvc := ucmocks.NewGrpcClient(t)
	mockSvc.EXPECT().Login(mock.Anything, "user@example.com", "password123").
		Return(nil, status.Error(codes.Unauthenticated, "bad credentials"))
	uc := setupUseCase(t, mockSvc)
	_, err := uc.Login(context.Background(), "user@example.com", "password123")
	assert.Error(t, err)
}
