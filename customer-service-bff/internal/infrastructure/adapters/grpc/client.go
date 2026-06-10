package grpcdapter

import (
	"context"
	"fmt"

	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/beabys/wms/customer-service-bff/internal/infrastructure/adapters/http/context"
)

// CustomerClient is a gRPC client for the customer-service.
// It wraps CustomerService and PermissionService RPCs.
type CustomerClient struct {
	customerService   customerv1.CustomerServiceClient
	permissionService customerv1.PermissionServiceClient
	conn              *grpc.ClientConn
}

// NewCustomerClient creates a new gRPC connection and client.
func NewCustomerClient(target string) (*CustomerClient, error) {
	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %w", err)
	}
	return &CustomerClient{
		customerService:   customerv1.NewCustomerServiceClient(conn),
		permissionService: customerv1.NewPermissionServiceClient(conn),
		conn:              conn,
	}, nil
}

// Close closes the gRPC connection.
func (c *CustomerClient) Close() error {
	return c.conn.Close()
}

// forwardJWTCtx extracts a JWT token from the context and attaches it
// as gRPC outgoing metadata if present.
func (c *CustomerClient) forwardJWTCtx(ctx context.Context) context.Context {
	token, _ := ctx.Value(httpctx.ContextKeyJWT).(string)
	if token != "" {
		return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
	}
	return ctx
}

// CustomerService RPCs

// RegisterCustomer registers a new customer.
func (c *CustomerClient) RegisterCustomer(ctx context.Context, req *customerv1.RegisterCustomerRequest) (*customerv1.RegisterCustomerResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.customerService.RegisterCustomer(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetCustomerByAdminID retrieves a customer by company admin user ID.
func (c *CustomerClient) GetCustomerByAdminID(ctx context.Context, adminUserID string) (*customerv1.GetCustomerByAdminIDResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.customerService.GetCustomerByAdminID(ctx, &customerv1.GetCustomerByAdminIDRequest{AdminUserId: adminUserID})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetCustomer retrieves a customer by ID.
func (c *CustomerClient) GetCustomer(ctx context.Context, id string) (*customerv1.GetCustomerResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.customerService.GetCustomer(ctx, &customerv1.GetCustomerRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateCustomer updates a customer's fields.
func (c *CustomerClient) UpdateCustomer(ctx context.Context, req *customerv1.UpdateCustomerRequest) (*customerv1.UpdateCustomerResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.customerService.UpdateCustomer(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListCustomers lists customers with optional filters and pagination.
func (c *CustomerClient) ListCustomers(ctx context.Context, req *customerv1.ListCustomersRequest) (*customerv1.ListCustomersResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.customerService.ListCustomers(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ApproveCustomer approves a customer.
func (c *CustomerClient) ApproveCustomer(ctx context.Context, id, approvedBy string) (*customerv1.ApproveCustomerResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.customerService.ApproveCustomer(ctx, &customerv1.ApproveCustomerRequest{
		Id:         id,
		ApprovedBy: approvedBy,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RejectCustomer rejects a customer.
func (c *CustomerClient) RejectCustomer(ctx context.Context, id, reason string) (*customerv1.RejectCustomerResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.customerService.RejectCustomer(ctx, &customerv1.RejectCustomerRequest{
		Id:     id,
		Reason: reason,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SuspendCustomer suspends a customer.
func (c *CustomerClient) SuspendCustomer(ctx context.Context, id, reason string) (*customerv1.SuspendCustomerResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.customerService.SuspendCustomer(ctx, &customerv1.SuspendCustomerRequest{
		Id:     id,
		Reason: reason,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RestoreCustomer restores a suspended customer.
func (c *CustomerClient) RestoreCustomer(ctx context.Context, id, restoredBy string) (*customerv1.RestoreCustomerResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.customerService.RestoreCustomer(ctx, &customerv1.RestoreCustomerRequest{
		Id:         id,
		RestoredBy: restoredBy,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListAuditLogs retrieves paginated audit logs for a customer.
func (c *CustomerClient) ListAuditLogs(ctx context.Context, req *customerv1.ListAuditLogsRequest) (*customerv1.ListAuditLogsResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.customerService.ListAuditLogs(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// PermissionService RPCs

// AssignCompanyRole assigns a company role to a user.
func (c *CustomerClient) AssignCompanyRole(ctx context.Context, req *customerv1.AssignCompanyRoleRequest) (*customerv1.AssignCompanyRoleResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.permissionService.AssignCompanyRole(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListCompanyRoles lists company roles.
func (c *CustomerClient) ListCompanyRoles(ctx context.Context, req *customerv1.ListCompanyRolesRequest) (*customerv1.ListCompanyRolesResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.permissionService.ListCompanyRoles(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetUserPermissions gets permissions for a user within a company.
func (c *CustomerClient) GetUserPermissions(ctx context.Context, req *customerv1.GetUserPermissionsRequest) (*customerv1.GetUserPermissionsResponse, error) {
	ctx = c.forwardJWTCtx(ctx)
	resp, err := c.permissionService.GetUserPermissions(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
