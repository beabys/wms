package usecase

import (
	"context"
	"log/slog"
	"testing"

	"github.com/beabys/wms/auth-service/internal/application/auth/command"
	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
	repomocks "github.com/beabys/wms/auth-service/mocks/application/auth/repository"
	"github.com/beabys/wms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupUserTest(t *testing.T) (*UserUseCase, *repomocks.UserRepository) {
	t.Helper()
	log := logger.NewSlogLogger(slog.LevelDebug)
	userRepo := repomocks.NewUserRepository(t)
	uc := NewUserUseCase(log, userRepo)
	return uc, userRepo
}

func TestUserUseCase_CreateUser(t *testing.T) {
	t.Run("successful creation returns user", func(t *testing.T) {
		uc, userRepo := setupUserTest(t)

		userRepo.On("GetByEmail", mock.Anything, "new@example.com").Return(nil, nil)
		userRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Email == "new@example.com" && u.Name == "New User" && u.Role == "admin"
		})).Return(nil)

		result, err := uc.CreateUser(context.Background(), command.CreateUserCommand{
			Email:    "new@example.com",
			Password: "password123",
			Name:     "New User",
			Role:     "admin",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, "new@example.com", result.Email)
		assert.Equal(t, "New User", result.Name)
		assert.Equal(t, "admin", result.Role)
		assert.True(t, result.Active)
	})

	t.Run("duplicate email returns error", func(t *testing.T) {
		uc, userRepo := setupUserTest(t)

		existingUser, _ := model.NewUser("dup@example.com", "password123", "First", "admin", nil)

		userRepo.On("GetByEmail", mock.Anything, "dup@example.com").Return(existingUser, nil)

		_, err := uc.CreateUser(context.Background(), command.CreateUserCommand{
			Email:    "dup@example.com",
			Password: "password456",
			Name:     "Second",
			Role:     "admin",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email already in use")
	})

	t.Run("empty email returns error", func(t *testing.T) {
		uc, _ := setupUserTest(t)

		_, err := uc.CreateUser(context.Background(), command.CreateUserCommand{
			Email:    "",
			Password: "password123",
			Name:     "Test",
			Role:     "admin",
		})
		assert.Error(t, err)
	})

	t.Run("empty password returns error", func(t *testing.T) {
		uc, _ := setupUserTest(t)

		_, err := uc.CreateUser(context.Background(), command.CreateUserCommand{
			Email:    "test@example.com",
			Password: "",
			Name:     "Test",
			Role:     "admin",
		})
		assert.Error(t, err)
	})

	t.Run("empty name returns error", func(t *testing.T) {
		uc, _ := setupUserTest(t)

		_, err := uc.CreateUser(context.Background(), command.CreateUserCommand{
			Email:    "test@example.com",
			Password: "password123",
			Name:     "",
			Role:     "admin",
		})
		assert.Error(t, err)
	})

	t.Run("empty role returns error", func(t *testing.T) {
		uc, _ := setupUserTest(t)

		_, err := uc.CreateUser(context.Background(), command.CreateUserCommand{
			Email:    "test@example.com",
			Password: "password123",
			Name:     "Test",
			Role:     "",
		})
		assert.Error(t, err)
	})
}

func TestUserUseCase_GetUser(t *testing.T) {
	t.Run("existing user returns user", func(t *testing.T) {
		uc, userRepo := setupUserTest(t)

		user, _ := model.NewUser("get@example.com", "password123", "Get User", "admin", nil)

		userRepo.On("GetByID", mock.Anything, user.ID).Return(user, nil)

		result, err := uc.GetUser(context.Background(), user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, "get@example.com", result.Email)
	})

	t.Run("non-existent user returns error", func(t *testing.T) {
		uc, userRepo := setupUserTest(t)

		userRepo.On("GetByID", mock.Anything, "non-existent-id").Return(nil, nil)

		_, err := uc.GetUser(context.Background(), "non-existent-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("empty ID returns error", func(t *testing.T) {
		uc, _ := setupUserTest(t)

		_, err := uc.GetUser(context.Background(), "")
		assert.Error(t, err)
	})
}

func TestUserUseCase_UpdateUser(t *testing.T) {
	t.Run("update name succeeds", func(t *testing.T) {
		uc, userRepo := setupUserTest(t)

		user, _ := model.NewUser("update@example.com", "password123", "Old Name", "admin", nil)

		userRepo.On("GetByID", mock.Anything, user.ID).Return(user, nil)
		userRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.ID == user.ID && u.Name == "New Name"
		})).Return(nil)

		newName := "New Name"
		result, err := uc.UpdateUser(context.Background(), command.UpdateUserCommand{
			UserID: user.ID,
			Name:   &newName,
		})
		require.NoError(t, err)
		assert.Equal(t, "New Name", result.Name)
	})

	t.Run("update role succeeds", func(t *testing.T) {
		uc, userRepo := setupUserTest(t)

		user, _ := model.NewUser("role@example.com", "password123", "Role User", "admin", nil)

		userRepo.On("GetByID", mock.Anything, user.ID).Return(user, nil)
		userRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.ID == user.ID && u.Role == "customer"
		})).Return(nil)

		newRole := "customer"
		result, err := uc.UpdateUser(context.Background(), command.UpdateUserCommand{
			UserID: user.ID,
			Role:   &newRole,
		})
		require.NoError(t, err)
		assert.Equal(t, "customer", result.Role)
	})

	t.Run("deactivate user succeeds", func(t *testing.T) {
		uc, userRepo := setupUserTest(t)

		user, _ := model.NewUser("deact@example.com", "password123", "Deact User", "admin", nil)

		userRepo.On("GetByID", mock.Anything, user.ID).Return(user, nil)
		userRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.ID == user.ID && !u.Active
		})).Return(nil)

		active := false
		result, err := uc.UpdateUser(context.Background(), command.UpdateUserCommand{
			UserID: user.ID,
			Active: &active,
		})
		require.NoError(t, err)
		assert.False(t, result.Active)
	})

	t.Run("non-existent user returns error", func(t *testing.T) {
		uc, userRepo := setupUserTest(t)

		userRepo.On("GetByID", mock.Anything, "non-existent").Return(nil, nil)

		_, err := uc.UpdateUser(context.Background(), command.UpdateUserCommand{
			UserID: "non-existent",
		})
		assert.Error(t, err)
	})

	t.Run("empty ID returns error", func(t *testing.T) {
		uc, _ := setupUserTest(t)

		_, err := uc.UpdateUser(context.Background(), command.UpdateUserCommand{
			UserID: "",
		})
		assert.Error(t, err)
	})
}

func TestUserUseCase_DeleteUser(t *testing.T) {
	t.Run("delete user succeeds", func(t *testing.T) {
		uc, userRepo := setupUserTest(t)

		user, _ := model.NewUser("delete@example.com", "password123", "Delete User", "admin", nil)

		userRepo.On("GetByID", mock.Anything, user.ID).Return(user, nil)
		userRepo.On("Deactivate", mock.Anything, user.ID).Return(nil)

		err := uc.DeleteUser(context.Background(), user.ID)
		require.NoError(t, err)

		userRepo.AssertCalled(t, "Deactivate", mock.Anything, user.ID)
	})

	t.Run("non-existent user returns error", func(t *testing.T) {
		uc, userRepo := setupUserTest(t)

		userRepo.On("GetByID", mock.Anything, "non-existent").Return(nil, nil)

		err := uc.DeleteUser(context.Background(), "non-existent")
		assert.Error(t, err)
	})

	t.Run("empty ID returns error", func(t *testing.T) {
		uc, _ := setupUserTest(t)

		err := uc.DeleteUser(context.Background(), "")
		assert.Error(t, err)
	})
}

func TestUserUseCase_AssignRole(t *testing.T) {
	t.Run("assign role succeeds", func(t *testing.T) {
		uc, userRepo := setupUserTest(t)

		user, _ := model.NewUser("role2@example.com", "password123", "Role User", "admin", nil)

		userRepo.On("GetByID", mock.Anything, user.ID).Return(user, nil)
		userRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.ID == user.ID && u.Role == "customer"
		})).Return(nil)

		err := uc.AssignRole(context.Background(), command.AssignRoleCommand{
			UserID: user.ID,
			Role:   "customer",
		})
		require.NoError(t, err)
	})

	t.Run("empty user ID returns error", func(t *testing.T) {
		uc, _ := setupUserTest(t)

		err := uc.AssignRole(context.Background(), command.AssignRoleCommand{
			UserID: "",
			Role:   "admin",
		})
		assert.Error(t, err)
	})

	t.Run("empty role returns error", func(t *testing.T) {
		uc, _ := setupUserTest(t)

		err := uc.AssignRole(context.Background(), command.AssignRoleCommand{
			UserID: "some-id",
			Role:   "",
		})
		assert.Error(t, err)
	})
}

func TestUserUseCase_ListUsers(t *testing.T) {
	t.Run("returns paginated results", func(t *testing.T) {
		uc, userRepo := setupUserTest(t)

		users := make([]*model.User, 5)
		for i := 0; i < 5; i++ {
			u, _ := model.NewUser(
				"user"+string(rune('0'+i))+"@example.com",
				"password123",
				"User"+string(rune('0'+i)),
				"admin",
				nil,
			)
			users[i] = u
		}

		userRepo.On("List", mock.Anything, mock.AnythingOfType("repository.UserFilter")).Return(users, 5, nil)

		result, err := uc.ListUsers(context.Background(), command.ListUsersQuery{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, 5, result.TotalCount)
		assert.Len(t, result.Users, 5)
	})

	t.Run("default pagination applies", func(t *testing.T) {
		uc, userRepo := setupUserTest(t)

		userRepo.On("List", mock.Anything, mock.AnythingOfType("repository.UserFilter")).Return([]*model.User{}, 0, nil)

		result, err := uc.ListUsers(context.Background(), command.ListUsersQuery{
			Page:     0,
			PageSize: 0,
		})
		require.NoError(t, err)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 20, result.PageSize)
	})
}

func TestUserUseCase_Stubs(t *testing.T) {
	t.Run("create role returns not implemented", func(t *testing.T) {
		uc, _ := setupUserTest(t)

		_, err := uc.CreateRole(context.Background(), command.CreateRoleCommand{
			Name: "test-role",
		})
		assert.Error(t, err)
	})

	t.Run("list roles returns not implemented", func(t *testing.T) {
		uc, _ := setupUserTest(t)

		_, err := uc.ListRoles(context.Background())
		assert.Error(t, err)
	})
}
