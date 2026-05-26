package repository

import (
	"context"

	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
)

// CustomerRepository defines the persistence contract for Customer aggregate.
// Defined on the consumer side (application layer).
type CustomerRepository interface {
	GetByID(ctx context.Context, id string) (*model.Customer, error)
	GetByEmail(ctx context.Context, email string) (*model.Customer, error)
	Save(ctx context.Context, customer *model.Customer) error
	List(ctx context.Context, status string, page, pageSize int32) ([]*model.Customer, int32, error)
	Update(ctx context.Context, customer *model.Customer) error
}

// InviteLinkRepository defines persistence for invitation links.
type InviteLinkRepository interface {
	Create(ctx context.Context, link *model.InviteLink) error
	GetByToken(ctx context.Context, token string) (*model.InviteLink, error)
	MarkUsed(ctx context.Context, id string) error
}
