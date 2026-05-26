package repository

import (
	"context"

	"github.com/beabys/wms/login-service/internal/domain/auth/model"
)

// UserRepository defines the persistence contract for users.
// Defined at the application layer (consumer side).
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email model.Email) (*model.User, error)
	Save(ctx context.Context, user *model.User) error
	List(ctx context.Context, limit, offset int32) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) error
	CreateRefreshToken(ctx context.Context, userID, token string) error
	ValidateRefreshToken(ctx context.Context, token string) (string, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}
