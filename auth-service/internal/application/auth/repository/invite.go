package repository

import (
	"context"
	"time"

	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
)

// InviteFilter contains optional filters for listing invites.
type InviteFilter struct {
	Status        *string
	Expired       *bool
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	Page          int
	PageSize      int
}

// InviteListResult contains the result of listing invites.
type InviteListResult struct {
	Invites    []*model.InviteToken
	TotalItems int
}

// InviteRepository defines the interface for invite token persistence.
type InviteRepository interface {
	// Create inserts a new invite token.
	Create(ctx context.Context, invite *model.InviteToken) error

	// GetByToken retrieves an invite token by its value. Returns nil, nil if not found.
	GetByToken(ctx context.Context, token string) (*model.InviteToken, error)

	// MarkUsed marks an invite token as used.
	MarkUsed(ctx context.Context, token string) error

	// Cancel sets invite status to cancelled (only if status=pending).
	Cancel(ctx context.Context, token string) error

	// List retrieves paginated invite tokens with optional filters.
	List(ctx context.Context, filter InviteFilter) (*InviteListResult, error)

	// ExistsPendingByEmail checks if there is a pending invite for the given email.
	ExistsPendingByEmail(ctx context.Context, email string) (bool, error)
}
