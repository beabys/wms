package usecase

import (
	"context"

	"github.com/beabys/wms/customer-service/internal/application/customer/command"
	"github.com/beabys/wms/customer-service/internal/application/customer/repository"
	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
)

// CustomerServiceHandler defines the application service contract for customer operations.
type CustomerServiceHandler interface {
	CreateCustomer(ctx context.Context, req *command.CreateCustomerCommand) (*model.Customer, error)
	InviteCustomer(ctx context.Context, req *command.InviteCustomerCommand) (*model.InviteLink, error)
	ApproveCustomer(ctx context.Context, req *command.ApproveCustomerCommand) (*model.Customer, error)
	SuspendCustomer(ctx context.Context, req *command.SuspendCustomerCommand) (*model.Customer, error)
	GetCustomer(ctx context.Context, req *command.GetCustomerQuery) (*model.Customer, error)
	ListCustomers(ctx context.Context, req *command.ListCustomersQuery) ([]*model.Customer, int32, error)
}

// CustomerService implements CustomerServiceHandler.
type CustomerService struct {
	customerRepo repository.CustomerRepository
	inviteRepo   repository.InviteLinkRepository
}
