package grpcdapter

import (
	"context"

	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"
)

// CustomerServiceAdapter wraps CustomerClient and maps gRPC errors to domain errors.
// It is a thin layer: proto in, proto out, no business logic, no model mapping.
type CustomerServiceAdapter struct {
	client *CustomerClient
}

// NewCustomerServiceAdapter creates a new adapter.
func NewCustomerServiceAdapter(client *CustomerClient) *CustomerServiceAdapter {
	return &CustomerServiceAdapter{client: client}
}

// RegisterCustomer registers a new customer.
func (a *CustomerServiceAdapter) RegisterCustomer(ctx context.Context, req *customerv1.RegisterCustomerRequest) (*customerv1.RegisterCustomerResponse, error) {
	resp, err := a.client.RegisterCustomer(ctx, req)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// GetCustomerByAdminID retrieves a customer by company admin user ID.
func (a *CustomerServiceAdapter) GetCustomerByAdminID(ctx context.Context, adminUserID string) (*customerv1.GetCustomerByAdminIDResponse, error) {
	resp, err := a.client.GetCustomerByAdminID(ctx, adminUserID)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// GetCustomer retrieves a customer by ID.
func (a *CustomerServiceAdapter) GetCustomer(ctx context.Context, id string) (*customerv1.GetCustomerResponse, error) {
	resp, err := a.client.GetCustomer(ctx, id)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// ListCustomers lists customers with optional filters.
func (a *CustomerServiceAdapter) ListCustomers(ctx context.Context, req *customerv1.ListCustomersRequest) (*customerv1.ListCustomersResponse, error) {
	resp, err := a.client.ListCustomers(ctx, req)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// UpdateCustomer updates a customer's fields.
func (a *CustomerServiceAdapter) UpdateCustomer(ctx context.Context, req *customerv1.UpdateCustomerRequest) (*customerv1.UpdateCustomerResponse, error) {
	resp, err := a.client.UpdateCustomer(ctx, req)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// ApproveCustomer approves a customer.
func (a *CustomerServiceAdapter) ApproveCustomer(ctx context.Context, id, approvedBy string) (*customerv1.ApproveCustomerResponse, error) {
	resp, err := a.client.ApproveCustomer(ctx, id, approvedBy)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// RejectCustomer rejects a customer.
func (a *CustomerServiceAdapter) RejectCustomer(ctx context.Context, id, reason string) (*customerv1.RejectCustomerResponse, error) {
	resp, err := a.client.RejectCustomer(ctx, id, reason)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// SuspendCustomer suspends a customer.
func (a *CustomerServiceAdapter) SuspendCustomer(ctx context.Context, id, reason string) (*customerv1.SuspendCustomerResponse, error) {
	resp, err := a.client.SuspendCustomer(ctx, id, reason)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// RestoreCustomer restores a suspended customer.
func (a *CustomerServiceAdapter) RestoreCustomer(ctx context.Context, id, restoredBy string) (*customerv1.RestoreCustomerResponse, error) {
	resp, err := a.client.RestoreCustomer(ctx, id, restoredBy)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// ListAuditLogs retrieves paginated audit logs for a customer.
func (a *CustomerServiceAdapter) ListAuditLogs(ctx context.Context, req *customerv1.ListAuditLogsRequest) (*customerv1.ListAuditLogsResponse, error) {
	resp, err := a.client.ListAuditLogs(ctx, req)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// AssignCompanyRole assigns a company role to a user.
func (a *CustomerServiceAdapter) AssignCompanyRole(ctx context.Context, req *customerv1.AssignCompanyRoleRequest) (*customerv1.AssignCompanyRoleResponse, error) {
	resp, err := a.client.AssignCompanyRole(ctx, req)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// ListCompanyRoles lists company roles.
func (a *CustomerServiceAdapter) ListCompanyRoles(ctx context.Context, req *customerv1.ListCompanyRolesRequest) (*customerv1.ListCompanyRolesResponse, error) {
	resp, err := a.client.ListCompanyRoles(ctx, req)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}

// GetUserPermissions gets permissions for a user within a company.
func (a *CustomerServiceAdapter) GetUserPermissions(ctx context.Context, req *customerv1.GetUserPermissionsRequest) (*customerv1.GetUserPermissionsResponse, error) {
	resp, err := a.client.GetUserPermissions(ctx, req)
	if err != nil {
		return nil, mapGrpcError(err)
	}
	return resp, nil
}
