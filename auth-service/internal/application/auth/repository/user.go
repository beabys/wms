package repository

import (
	"context"

	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
)

// UserFilter specifies filtering and pagination for listing users.
type UserFilter struct {
	CustomerID *string
	Role       *string
	Active     *bool
	Page       int
	PageSize   int
}

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	// GetByID retrieves a user by their ID. Returns nil, nil if not found.
	GetByID(ctx context.Context, id string) (*model.User, error)

	// GetByEmail retrieves a user by their email. Returns nil, nil if not found.
	GetByEmail(ctx context.Context, email string) (*model.User, error)

	// Create inserts a new user.
	Create(ctx context.Context, user *model.User) error

	// Update updates an existing user.
	Update(ctx context.Context, user *model.User) error

	// List retrieves users based on filter with pagination.
	List(ctx context.Context, filter UserFilter) ([]*model.User, int, error)

	// Deactivate sets a user's active flag to false.
	Deactivate(ctx context.Context, id string) error
}
