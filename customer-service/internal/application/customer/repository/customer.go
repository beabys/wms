package repository

import (
	"context"

	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
)

// CustomerRepository defines the interface for customer persistence.
type CustomerRepository interface {
	// Create inserts a new customer.
	Create(ctx context.Context, customer *model.Customer) error

	// GetByID retrieves a customer by ID. Returns nil, nil if not found.
	GetByID(ctx context.Context, id string) (*model.Customer, error)

	// GetByEmail retrieves a customer by email. Returns nil, nil if not found.
	GetByEmail(ctx context.Context, email string) (*model.Customer, error)

	// Update updates an existing customer.
	Update(ctx context.Context, customer *model.Customer) error

	// List retrieves customers with filtering and pagination.
	List(ctx context.Context, page, pageSize int, status string) ([]*model.Customer, int, error)

	// UpdateStatus changes a customer's status.
	UpdateStatus(ctx context.Context, id string, status model.CustomerStatus, reason string) error

	// FindByAdminID retrieves a customer by company_admin_id. Returns nil, nil if not found.
	FindByAdminID(ctx context.Context, adminUserID string) (*model.Customer, error)
}
