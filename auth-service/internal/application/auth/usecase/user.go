package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/beabys/wms/auth-service/internal/application/auth/command"
	"github.com/beabys/wms/auth-service/internal/application/auth/repository"
	"github.com/beabys/wms/auth-service/internal/application/auth/validator"
	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
	"github.com/beabys/wms/pkg/logger"
)

// UserUseCase handles user management operations.
type UserUseCase struct {
	logger   logger.Logger
	userRepo repository.UserRepository
}

// NewUserUseCase creates a new UserUseCase.
func NewUserUseCase(log logger.Logger, userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		logger:   log,
		userRepo: userRepo,
	}
}

// CreateUser creates a new user in the system.
func (uc *UserUseCase) CreateUser(ctx context.Context, cmd command.CreateUserCommand) (*command.UserResult, error) {
	if err := validator.Email(cmd.Email); err != nil {
		return nil, err
	}
	if err := validator.Password(cmd.Password); err != nil {
		return nil, err
	}
	if err := validator.NotEmpty(cmd.Name, "name"); err != nil {
		return nil, err
	}
	if err := validator.NotEmpty(cmd.Role, "role"); err != nil {
		return nil, err
	}

	// Check email uniqueness
	existing, err := uc.userRepo.GetByEmail(ctx, cmd.Email)
	if err != nil {
		uc.logger.Error("failed to check email uniqueness", err)
		return nil, fmt.Errorf("failed to create user")
	}
	if existing != nil {
		return nil, fmt.Errorf("email already in use")
	}

	user, err := model.NewUser(cmd.Email, cmd.Password, cmd.Name, cmd.Role, cmd.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("invalid user data: %w", err)
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		uc.logger.Error("failed to create user", err)
		return nil, fmt.Errorf("failed to create user")
	}

	return userToResult(user), nil
}

// GetUser retrieves a user by ID.
func (uc *UserUseCase) GetUser(ctx context.Context, id string) (*command.UserResult, error) {
	if err := validator.NotEmpty(id, "user ID"); err != nil {
		return nil, err
	}

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("failed to get user", err,
			logger.LogField{Key: "user_id", Value: id},
		)
		return nil, fmt.Errorf("failed to get user")
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return userToResult(user), nil
}

// UpdateUser updates an existing user's fields.
func (uc *UserUseCase) UpdateUser(ctx context.Context, cmd command.UpdateUserCommand) (*command.UserResult, error) {
	if err := validator.NotEmpty(cmd.UserID, "user ID"); err != nil {
		return nil, err
	}

	user, err := uc.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		uc.logger.Error("failed to get user for update", err)
		return nil, fmt.Errorf("failed to update user")
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if cmd.Name != nil {
		user.Name = *cmd.Name
	}
	if cmd.Role != nil {
		user.Role = *cmd.Role
	}
	if cmd.Active != nil {
		user.Active = *cmd.Active
	}
	user.UpdatedAt = now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Error("failed to update user", err)
		return nil, fmt.Errorf("failed to update user")
	}

	return userToResult(user), nil
}

// ListUsers retrieves users with optional filtering and pagination.
func (uc *UserUseCase) ListUsers(ctx context.Context, query command.ListUsersQuery) (*command.ListUsersResult, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	filter := repository.UserFilter{
		CustomerID: query.CustomerID,
		Role:       query.Role,
		Active:     query.Active,
		Page:       query.Page,
		PageSize:   query.PageSize,
	}

	users, total, err := uc.userRepo.List(ctx, filter)
	if err != nil {
		uc.logger.Error("failed to list users", err)
		return nil, fmt.Errorf("failed to list users")
	}

	results := make([]command.UserResult, 0, len(users))
	for _, u := range users {
		results = append(results, *userToResult(u))
	}

	return &command.ListUsersResult{
		Users:      results,
		TotalCount: total,
		Page:       query.Page,
		PageSize:   query.PageSize,
	}, nil
}

// DeleteUser deactivates a user (soft delete).
func (uc *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	if err := validator.NotEmpty(id, "user ID"); err != nil {
		return err
	}

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("failed to get user for delete", err)
		return fmt.Errorf("failed to delete user")
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if err := uc.userRepo.Deactivate(ctx, id); err != nil {
		uc.logger.Error("failed to deactivate user", err)
		return fmt.Errorf("failed to delete user")
	}

	return nil
}

// AssignRole assigns a role to a user.
func (uc *UserUseCase) AssignRole(ctx context.Context, cmd command.AssignRoleCommand) error {
	if err := validator.NotEmpty(cmd.UserID, "user ID"); err != nil {
		return err
	}
	if err := validator.NotEmpty(cmd.Role, "role"); err != nil {
		return err
	}

	user, err := uc.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		uc.logger.Error("failed to get user for role assignment", err)
		return fmt.Errorf("failed to assign role")
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	user.Role = cmd.Role
	user.UpdatedAt = now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Error("failed to update user role", err)
		return fmt.Errorf("failed to assign role")
	}

	return nil
}

// CreateRole is a stub for creating a role (full implementation in future).
func (uc *UserUseCase) CreateRole(ctx context.Context, cmd command.CreateRoleCommand) (*model.Role, error) {
	return nil, fmt.Errorf("not implemented")
}

// ListRoles is a stub for listing roles (full implementation in future).
func (uc *UserUseCase) ListRoles(ctx context.Context) ([]*model.Role, error) {
	return nil, fmt.Errorf("not implemented")
}

func userToResult(u *model.User) *command.UserResult {
	return &command.UserResult{
		ID:         u.ID,
		Email:      u.Email,
		Name:       u.Name,
		Role:       u.Role,
		CustomerID: u.CustomerID,
		Active:     u.Active,
		CreatedAt:  u.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:  u.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func now() time.Time {
	return time.Now()
}
