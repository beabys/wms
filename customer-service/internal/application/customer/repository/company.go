package repository

import (
	"context"

	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
)

// CompanyRepository defines the interface for company persistence.
type CompanyRepository interface {
	// Create inserts a new company.
	Create(ctx context.Context, company *model.Company) error

	// GetByID retrieves a company by ID. Returns nil, nil if not found.
	GetByID(ctx context.Context, id string) (*model.Company, error)

	// GetByCustomerID retrieves a company by customer ID. Returns nil, nil if not found.
	GetByCustomerID(ctx context.Context, customerID string) (*model.Company, error)
}

// CompanyRoleRepository defines the interface for company role persistence.
type CompanyRoleRepository interface {
	// Create inserts a new company role.
	Create(ctx context.Context, role *model.CompanyRole) error

	// ListByCustomer retrieves all roles for a customer.
	ListByCustomer(ctx context.Context, customerID string) ([]*model.CompanyRole, error)

	// GetUserPermissions retrieves aggregated permissions for a user in a customer.
	GetUserPermissions(ctx context.Context, userID, customerID string) ([]string, error)

	// AssignUserRole assigns a role to a user within a customer.
	AssignUserRole(ctx context.Context, userID, roleID, customerID, assignedBy string) error
}
