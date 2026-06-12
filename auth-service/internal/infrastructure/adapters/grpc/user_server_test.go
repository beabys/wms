package grpcdapter

import (
	"context"
	"testing"

	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
	repomocks "github.com/beabys/wms/auth-service/mocks/application/auth/repository"
	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUserServer_CreateUser(t *testing.T) {
	t.Run("create user succeeds", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewUserServiceClient(conn)

		userRepo.On("GetByEmail", mock.Anything, "newuser@example.com").Return(nil, nil).Once()
		userRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Email == "newuser@example.com" && u.Name == "New User" && u.Role == "admin"
		})).Return(nil).Once()

		resp, err := client.CreateUser(context.Background(), &authv1.CreateUserRequest{
			Email:    "newuser@example.com",
			Password: "password123",
			Name:     "New User",
			Role:     "admin",
		})
		require.NoError(t, err)
		require.NotNil(t, resp.GetUser())
		assert.NotEmpty(t, resp.GetUser().GetId())
		assert.Equal(t, "newuser@example.com", resp.GetUser().GetEmail())
		assert.Equal(t, "admin", resp.GetUser().GetRole())
		assert.True(t, resp.GetUser().GetActive())
	})

	t.Run("duplicate email returns error", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewUserServiceClient(conn)

		existingUser, _ := model.NewUser("newuser@example.com", "password123", "Existing", "admin", nil)

		userRepo.On("GetByEmail", mock.Anything, "newuser@example.com").Return(existingUser, nil).Once()

		_, err := client.CreateUser(context.Background(), &authv1.CreateUserRequest{
			Email:    "newuser@example.com",
			Password: "password123",
			Name:     "Duplicate",
			Role:     "admin",
		})
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("missing fields return error", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewUserServiceClient(conn)

		_, err := client.CreateUser(context.Background(), &authv1.CreateUserRequest{
			Email:    "",
			Password: "",
			Name:     "",
			Role:     "",
		})
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})
}

func TestUserServer_GetUser(t *testing.T) {
	t.Run("get existing user succeeds", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewUserServiceClient(conn)

		// Create
		userRepo.On("GetByEmail", mock.Anything, "getuser@example.com").Return(nil, nil).Once()
		userRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Email == "getuser@example.com"
		})).Return(nil).Once()

		createResp, err := client.CreateUser(context.Background(), &authv1.CreateUserRequest{
			Email:    "getuser@example.com",
			Password: "password123",
			Name:     "Get User",
			Role:     "admin",
		})
		require.NoError(t, err)
		userID := createResp.GetUser().GetId()

		// Get
		userRepo.On("GetByID", mock.Anything, userID).Return(&model.User{
			ID:    userID,
			Email: "getuser@example.com",
			Name:  "Get User",
			Role:  "admin",
			Active: true,
		}, nil).Once()

		resp, err := client.GetUser(context.Background(), &authv1.GetUserRequest{
			Id: userID,
		})
		require.NoError(t, err)
		assert.Equal(t, "getuser@example.com", resp.GetUser().GetEmail())
		assert.Equal(t, "Get User", resp.GetUser().GetName())
	})

	t.Run("get non-existent user returns error", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewUserServiceClient(conn)

		userRepo.On("GetByID", mock.Anything, "non-existent-id").Return(nil, nil).Once()

		_, err := client.GetUser(context.Background(), &authv1.GetUserRequest{
			Id: "non-existent-id",
		})
		require.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
	})
}

func TestUserServer_UpdateUser(t *testing.T) {
	t.Run("update user name succeeds", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewUserServiceClient(conn)

		// Create
		userRepo.On("GetByEmail", mock.Anything, "updateuser@example.com").Return(nil, nil).Once()
		userRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Email == "updateuser@example.com"
		})).Return(nil).Once()

		createResp, err := client.CreateUser(context.Background(), &authv1.CreateUserRequest{
			Email:    "updateuser@example.com",
			Password: "password123",
			Name:     "Old Name",
			Role:     "admin",
		})
		require.NoError(t, err)
		userID := createResp.GetUser().GetId()

		// Update
		userRepo.On("GetByID", mock.Anything, userID).Return(&model.User{
			ID:    userID,
			Email: "updateuser@example.com",
			Name:  "Old Name",
			Role:  "admin",
			Active: true,
		}, nil).Once()
		userRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.ID == userID && u.Name == "New Name"
		})).Return(nil).Once()

		resp, err := client.UpdateUser(context.Background(), &authv1.UpdateUserRequest{
			Id:   userID,
			Name: "New Name",
		})
		require.NoError(t, err)
		assert.Equal(t, "New Name", resp.GetUser().GetName())
	})

	t.Run("update non-existent user returns error", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewUserServiceClient(conn)

		userRepo.On("GetByID", mock.Anything, "non-existent").Return(nil, nil).Once()

		_, err := client.UpdateUser(context.Background(), &authv1.UpdateUserRequest{
			Id: "non-existent",
		})
		require.Error(t, err)
	})
}

func TestUserServer_DeleteUser(t *testing.T) {
	t.Run("delete user succeeds", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewUserServiceClient(conn)

		// Create
		userRepo.On("GetByEmail", mock.Anything, "deleteuser@example.com").Return(nil, nil).Once()
		userRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Email == "deleteuser@example.com"
		})).Return(nil).Once()

		createResp, err := client.CreateUser(context.Background(), &authv1.CreateUserRequest{
			Email:    "deleteuser@example.com",
			Password: "password123",
			Name:     "Delete User",
			Role:     "admin",
		})
		require.NoError(t, err)
		userID := createResp.GetUser().GetId()

		// Delete
		userRepo.On("GetByID", mock.Anything, userID).Return(&model.User{
			ID:    userID,
			Email: "deleteuser@example.com",
			Name:  "Delete User",
			Role:  "admin",
			Active: true,
		}, nil).Once()
		userRepo.On("Deactivate", mock.Anything, userID).Return(nil).Once()

		resp, err := client.DeleteUser(context.Background(), &authv1.DeleteUserRequest{
			Id: userID,
		})
		require.NoError(t, err)
		assert.True(t, resp.GetSuccess())
	})

	t.Run("delete non-existent user returns error", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewUserServiceClient(conn)

		userRepo.On("GetByID", mock.Anything, "non-existent").Return(nil, nil).Once()

		_, err := client.DeleteUser(context.Background(), &authv1.DeleteUserRequest{
			Id: "non-existent",
		})
		require.Error(t, err)
	})
}

func TestUserServer_AssignRole(t *testing.T) {
	t.Run("assign role succeeds", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewUserServiceClient(conn)

		// Create
		userRepo.On("GetByEmail", mock.Anything, "assignrole@example.com").Return(nil, nil).Once()
		userRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Email == "assignrole@example.com"
		})).Return(nil).Once()

		createResp, err := client.CreateUser(context.Background(), &authv1.CreateUserRequest{
			Email:    "assignrole@example.com",
			Password: "password123",
			Name:     "Assign Role",
			Role:     "admin",
		})
		require.NoError(t, err)
		userID := createResp.GetUser().GetId()

		// AssignRole
		userRepo.On("GetByID", mock.Anything, userID).Return(&model.User{
			ID:    userID,
			Email: "assignrole@example.com",
			Name:  "Assign Role",
			Role:  "admin",
			Active: true,
		}, nil).Once()
		userRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.ID == userID && u.Role == "customer"
		})).Return(nil).Once()

		resp, err := client.AssignRole(context.Background(), &authv1.AssignRoleRequest{
			UserId: userID,
			Role:   "customer",
		})
		require.NoError(t, err)
		assert.True(t, resp.GetSuccess())
	})

	t.Run("assign role to non-existent user returns error", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewUserServiceClient(conn)

		userRepo.On("GetByID", mock.Anything, "non-existent").Return(nil, nil).Once()

		_, err := client.AssignRole(context.Background(), &authv1.AssignRoleRequest{
			UserId: "non-existent",
			Role:   "admin",
		})
		require.Error(t, err)
	})
}

func TestUserServer_ListUsers(t *testing.T) {
	t.Run("list users returns paginated results", func(t *testing.T) {
		userRepo := repomocks.NewUserRepository(t)
		authRepo := repomocks.NewAuthRepository(t)
		inviteRepo := repomocks.NewInviteRepository(t)

		conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
		client := authv1.NewUserServiceClient(conn)

		user, _ := model.NewUser("listuser@example.com", "password123", "List User", "admin", nil)

		// Create
		userRepo.On("GetByEmail", mock.Anything, "listuser@example.com").Return(nil, nil).Once()
		userRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Email == "listuser@example.com"
		})).Return(nil).Once()

		_, err := client.CreateUser(context.Background(), &authv1.CreateUserRequest{
			Email:    "listuser@example.com",
			Password: "password123",
			Name:     "List User",
			Role:     "admin",
		})
		require.NoError(t, err)

		// List
		userRepo.On("List", mock.Anything, mock.AnythingOfType("repository.UserFilter")).Return([]*model.User{user}, 1, nil).Once()

		resp, err := client.ListUsers(context.Background(), &authv1.ListUsersRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.NotNil(t, resp.GetUsers())
		assert.GreaterOrEqual(t, resp.GetPagination().GetTotal(), int32(1))
	})
}

func TestUserServer_Unimplemented(t *testing.T) {
	userRepo := repomocks.NewUserRepository(t)
	authRepo := repomocks.NewAuthRepository(t)
	inviteRepo := repomocks.NewInviteRepository(t)

	conn, _ := setupGRPCServer(t, serverDeps{userRepo, authRepo, inviteRepo})
	client := authv1.NewUserServiceClient(conn)

	t.Run("list roles returns unimplemented", func(t *testing.T) {
		_, err := client.ListRoles(context.Background(), &authv1.ListRolesRequest{})
		require.Error(t, err)
		assert.Equal(t, codes.Unimplemented, status.Code(err))
	})

	t.Run("create role returns unimplemented", func(t *testing.T) {
		_, err := client.CreateRole(context.Background(), &authv1.CreateRoleRequest{
			Name: "test",
		})
		require.Error(t, err)
		assert.Equal(t, codes.Unimplemented, status.Code(err))
	})
}
