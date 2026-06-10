package usecase

import (
	"context"
	"fmt"

	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"

	"github.com/beabys/wms/customer-service-bff/internal/application/customer/transformer"
	"github.com/beabys/wms/customer-service-bff/internal/application/customer/validator"
	grpcdapter "github.com/beabys/wms/customer-service-bff/internal/infrastructure/adapters/grpc"
	"github.com/beabys/wms/customer-service-bff/internal/domain/model"
	"github.com/beabys/wms/pkg/logger"
)

// GrpcClient defines what the use case needs from the gRPC layer.
// Interface defined by consumer (usecase).
type GrpcClient interface {
	RegisterCustomer(ctx context.Context, req *customerv1.RegisterCustomerRequest) (*customerv1.RegisterCustomerResponse, error)
	GetCustomerByAdminID(ctx context.Context, adminUserID string) (*customerv1.GetCustomerByAdminIDResponse, error)
	GetCustomer(ctx context.Context, id string) (*customerv1.GetCustomerResponse, error)
	ListCustomers(ctx context.Context, req *customerv1.ListCustomersRequest) (*customerv1.ListCustomersResponse, error)
	UpdateCustomer(ctx context.Context, req *customerv1.UpdateCustomerRequest) (*customerv1.UpdateCustomerResponse, error)
	ApproveCustomer(ctx context.Context, id, approvedBy string) (*customerv1.ApproveCustomerResponse, error)
	RejectCustomer(ctx context.Context, id, reason string) (*customerv1.RejectCustomerResponse, error)
	SuspendCustomer(ctx context.Context, id, reason string) (*customerv1.SuspendCustomerResponse, error)
	RestoreCustomer(ctx context.Context, id, restoredBy string) (*customerv1.RestoreCustomerResponse, error)
	AssignCompanyRole(ctx context.Context, req *customerv1.AssignCompanyRoleRequest) (*customerv1.AssignCompanyRoleResponse, error)
	ListCompanyRoles(ctx context.Context, req *customerv1.ListCompanyRolesRequest) (*customerv1.ListCompanyRolesResponse, error)
	GetUserPermissions(ctx context.Context, req *customerv1.GetUserPermissionsRequest) (*customerv1.GetUserPermissionsResponse, error)
	ListAuditLogs(ctx context.Context, req *customerv1.ListAuditLogsRequest) (*customerv1.ListAuditLogsResponse, error)
}

// CustomerUseCase implements business logic for customer operations.
type CustomerUseCase struct {
	grpc GrpcClient
	log  logger.Logger
}

// New creates a new CustomerUseCase.
func New(log logger.Logger, grpc GrpcClient) *CustomerUseCase {
	return &CustomerUseCase{grpc: grpc, log: log}
}

// RegisterCustomer registers a new customer.
func (uc *CustomerUseCase) RegisterCustomer(ctx context.Context, req *model.RegisterCustomerRequest) (*model.RegisterCustomerResponse, error) {
	if err := validator.ValidateRegisterCustomer(req.Token, req.CompanyName, req.Email, req.Password); err != nil {
		return nil, fmt.Errorf("%w: %w", grpcdapter.ErrInvalidArgument, err)
	}

	protoReq := &customerv1.RegisterCustomerRequest{
		Token:       req.Token,
		CompanyName: req.CompanyName,
		Email:       req.Email,
		Password:    req.Password,
		Phone:       req.Phone,
		VatNumber:   req.VatNumber,
		Address:     req.Address,
		City:        req.City,
		PostalCode:  req.PostalCode,
		Country:     req.Country,
	}

	resp, err := uc.grpc.RegisterCustomer(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.RegisterCustomerResponse{
		Customer:    transformer.CustomerFromGrpc(resp.GetCustomer()),
		AccessToken: resp.GetAccessToken(),
	}, nil
}

// GetMyCustomer retrieves the customer profile for the authenticated user.
func (uc *CustomerUseCase) GetMyCustomer(ctx context.Context, adminUserID string) (*model.CustomerResponse, error) {
	resp, err := uc.grpc.GetCustomerByAdminID(ctx, adminUserID)
	if err != nil {
		return nil, err
	}
	return transformer.CustomerFromGrpc(resp.GetCustomer()), nil
}

// GetCustomer retrieves a customer by ID.
func (uc *CustomerUseCase) GetCustomer(ctx context.Context, id string) (*model.CustomerResponse, error) {
	resp, err := uc.grpc.GetCustomer(ctx, id)
	if err != nil {
		return nil, err
	}
	return transformer.CustomerFromGrpc(resp.GetCustomer()), nil
}

// ListCustomers lists customers with optional filters.
func (uc *CustomerUseCase) ListCustomers(ctx context.Context, status *string, page, pageSize *int) (*model.CustomerListResponse, error) {
	protoReq := &customerv1.ListCustomersRequest{}
	if status != nil {
		protoReq.Status = *status
	}
	if page != nil {
		protoReq.Page = int32(*page)
	}
	if pageSize != nil {
		protoReq.PageSize = int32(*pageSize)
	}

	resp, err := uc.grpc.ListCustomers(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.CustomerListResponse{
		Customers:  transformer.CustomerListFromGrpc(resp.GetCustomers()),
		Pagination: transformer.PaginationFromGrpc(resp.GetPagination()),
	}, nil
}

// UpdateCustomer updates a customer's fields.
func (uc *CustomerUseCase) UpdateCustomer(ctx context.Context, id string, req *model.UpdateCustomerRequest) (*model.CustomerResponse, error) {
	protoReq := &customerv1.UpdateCustomerRequest{
		Id:         id,
		Phone:      req.Phone,
		Address:    req.Address,
		City:       req.City,
		PostalCode: req.PostalCode,
		Country:    req.Country,
	}

	resp, err := uc.grpc.UpdateCustomer(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return transformer.CustomerFromGrpc(resp.GetCustomer()), nil
}

// ApproveCustomer approves a customer.
func (uc *CustomerUseCase) ApproveCustomer(ctx context.Context, id, approvedBy string) (*model.CustomerResponse, error) {
	resp, err := uc.grpc.ApproveCustomer(ctx, id, approvedBy)
	if err != nil {
		return nil, err
	}
	return transformer.CustomerFromGrpc(resp.GetCustomer()), nil
}

// RejectCustomer rejects a customer.
func (uc *CustomerUseCase) RejectCustomer(ctx context.Context, id, reason string) error {
	_, err := uc.grpc.RejectCustomer(ctx, id, reason)
	return err
}

// SuspendCustomer suspends a customer.
func (uc *CustomerUseCase) SuspendCustomer(ctx context.Context, id, reason string) error {
	_, err := uc.grpc.SuspendCustomer(ctx, id, reason)
	return err
}

// RestoreCustomer restores a suspended customer.
func (uc *CustomerUseCase) RestoreCustomer(ctx context.Context, id, restoredBy string) (*model.CustomerResponse, error) {
	resp, err := uc.grpc.RestoreCustomer(ctx, id, restoredBy)
	if err != nil {
		return nil, err
	}
	return transformer.CustomerFromGrpc(resp.GetCustomer()), nil
}

// ListAuditLogs retrieves paginated audit logs for a customer.
func (uc *CustomerUseCase) ListAuditLogs(ctx context.Context, customerID string, page, pageSize int) (*model.AuditLogListResponse, error) {
	protoReq := &customerv1.ListAuditLogsRequest{
		CustomerId: customerID,
		Page:       int32(page),
		PageSize:   int32(pageSize),
	}

	resp, err := uc.grpc.ListAuditLogs(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	entries := make([]model.AuditEntry, 0, len(resp.GetEntries()))
	for _, e := range resp.GetEntries() {
		entries = append(entries, transformer.AuditEntryFromGrpc(e))
	}

	return &model.AuditLogListResponse{
		Entries:    entries,
		Pagination: transformer.PaginationFromGrpc(resp.GetPagination()),
	}, nil
}

// AssignCompanyRole assigns a company role to a user.
func (uc *CustomerUseCase) AssignCompanyRole(ctx context.Context, req *model.AssignCompanyRoleRequest) error {
	if err := validator.ValidateAssignCompanyRole(req.CustomerID, req.UserID, req.RoleName); err != nil {
		return fmt.Errorf("%w: %w", grpcdapter.ErrInvalidArgument, err)
	}

	protoReq := &customerv1.AssignCompanyRoleRequest{
		CustomerId:  req.CustomerID,
		UserId:      req.UserID,
		RoleName:    req.RoleName,
		Permissions: req.Permissions,
		AssignedBy:  req.AssignedBy,
	}

	_, err := uc.grpc.AssignCompanyRole(ctx, protoReq)
	return err
}

// ListCompanyRoles lists company roles.
func (uc *CustomerUseCase) ListCompanyRoles(ctx context.Context, customerID string) (*model.CompanyRoleListResponse, error) {
	resp, err := uc.grpc.ListCompanyRoles(ctx, &customerv1.ListCompanyRolesRequest{
		CustomerId: customerID,
	})
	if err != nil {
		return nil, err
	}

	roles := make([]model.CompanyRoleResponse, 0, len(resp.GetRoles()))
	for _, r := range resp.GetRoles() {
		roles = append(roles, transformer.CompanyRoleFromGrpc(r))
	}
	return &model.CompanyRoleListResponse{Roles: roles}, nil
}

// GetUserPermissions gets permissions for a user within a company.
func (uc *CustomerUseCase) GetUserPermissions(ctx context.Context, userID, customerID string) (*model.PermissionsResponse, error) {
	resp, err := uc.grpc.GetUserPermissions(ctx, &customerv1.GetUserPermissionsRequest{
		UserId:     userID,
		CustomerId: customerID,
	})
	if err != nil {
		return nil, err
	}
	return &model.PermissionsResponse{
		Permissions: resp.GetPermissions(),
	}, nil
}
