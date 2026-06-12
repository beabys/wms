package grpcdapter

import (
	"context"
	"net"
	"testing"

	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"
	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// testCustomerServer implements customerv1.CustomerServiceServer for testing.
type testCustomerServer struct {
	customerv1.UnimplementedCustomerServiceServer
	registerCustomerFn func(ctx context.Context, req *customerv1.RegisterCustomerRequest) (*customerv1.RegisterCustomerResponse, error)
	getCustomerFn      func(ctx context.Context, req *customerv1.GetCustomerRequest) (*customerv1.GetCustomerResponse, error)
	updateCustomerFn   func(ctx context.Context, req *customerv1.UpdateCustomerRequest) (*customerv1.UpdateCustomerResponse, error)
	listCustomersFn    func(ctx context.Context, req *customerv1.ListCustomersRequest) (*customerv1.ListCustomersResponse, error)
	approveCustomerFn  func(ctx context.Context, req *customerv1.ApproveCustomerRequest) (*customerv1.ApproveCustomerResponse, error)
	rejectCustomerFn   func(ctx context.Context, req *customerv1.RejectCustomerRequest) (*customerv1.RejectCustomerResponse, error)
	suspendCustomerFn  func(ctx context.Context, req *customerv1.SuspendCustomerRequest) (*customerv1.SuspendCustomerResponse, error)
}

func (s *testCustomerServer) RegisterCustomer(ctx context.Context, req *customerv1.RegisterCustomerRequest) (*customerv1.RegisterCustomerResponse, error) {
	if s.registerCustomerFn != nil {
		return s.registerCustomerFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *testCustomerServer) GetCustomer(ctx context.Context, req *customerv1.GetCustomerRequest) (*customerv1.GetCustomerResponse, error) {
	if s.getCustomerFn != nil {
		return s.getCustomerFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *testCustomerServer) UpdateCustomer(ctx context.Context, req *customerv1.UpdateCustomerRequest) (*customerv1.UpdateCustomerResponse, error) {
	if s.updateCustomerFn != nil {
		return s.updateCustomerFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *testCustomerServer) ListCustomers(ctx context.Context, req *customerv1.ListCustomersRequest) (*customerv1.ListCustomersResponse, error) {
	if s.listCustomersFn != nil {
		return s.listCustomersFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *testCustomerServer) ApproveCustomer(ctx context.Context, req *customerv1.ApproveCustomerRequest) (*customerv1.ApproveCustomerResponse, error) {
	if s.approveCustomerFn != nil {
		return s.approveCustomerFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *testCustomerServer) RejectCustomer(ctx context.Context, req *customerv1.RejectCustomerRequest) (*customerv1.RejectCustomerResponse, error) {
	if s.rejectCustomerFn != nil {
		return s.rejectCustomerFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *testCustomerServer) SuspendCustomer(ctx context.Context, req *customerv1.SuspendCustomerRequest) (*customerv1.SuspendCustomerResponse, error) {
	if s.suspendCustomerFn != nil {
		return s.suspendCustomerFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// testPermissionServer implements customerv1.PermissionServiceServer for testing.
type testPermissionServer struct {
	customerv1.UnimplementedPermissionServiceServer
	assignCompanyRoleFn  func(ctx context.Context, req *customerv1.AssignCompanyRoleRequest) (*customerv1.AssignCompanyRoleResponse, error)
	listCompanyRolesFn   func(ctx context.Context, req *customerv1.ListCompanyRolesRequest) (*customerv1.ListCompanyRolesResponse, error)
	getUserPermissionsFn func(ctx context.Context, req *customerv1.GetUserPermissionsRequest) (*customerv1.GetUserPermissionsResponse, error)
}

func (s *testPermissionServer) AssignCompanyRole(ctx context.Context, req *customerv1.AssignCompanyRoleRequest) (*customerv1.AssignCompanyRoleResponse, error) {
	if s.assignCompanyRoleFn != nil {
		return s.assignCompanyRoleFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *testPermissionServer) ListCompanyRoles(ctx context.Context, req *customerv1.ListCompanyRolesRequest) (*customerv1.ListCompanyRolesResponse, error) {
	if s.listCompanyRolesFn != nil {
		return s.listCompanyRolesFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *testPermissionServer) GetUserPermissions(ctx context.Context, req *customerv1.GetUserPermissionsRequest) (*customerv1.GetUserPermissionsResponse, error) {
	if s.getUserPermissionsFn != nil {
		return s.getUserPermissionsFn(ctx, req)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// setupTestServer starts a local gRPC server and returns a client connected to it.
func setupTestServer(t *testing.T, custSrv *testCustomerServer, permSrv *testPermissionServer) (*CustomerClient, func()) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	srv := grpc.NewServer()
	if custSrv != nil {
		customerv1.RegisterCustomerServiceServer(srv, custSrv)
	}
	if permSrv != nil {
		customerv1.RegisterPermissionServiceServer(srv, permSrv)
	}

	go func() {
		srv.Serve(listener) //nolint:errcheck
	}()

	client, err := NewCustomerClient(listener.Addr().String())
	require.NoError(t, err)

	cleanup := func() {
		client.Close()
		srv.GracefulStop()
	}

	return client, cleanup
}

func TestCustomerServiceAdapterRegisterCustomer(t *testing.T) {
	custSrv := &testCustomerServer{
		registerCustomerFn: func(ctx context.Context, req *customerv1.RegisterCustomerRequest) (*customerv1.RegisterCustomerResponse, error) {
			return &customerv1.RegisterCustomerResponse{
				Customer: &customerv1.Customer{
					Id:          "new-cust-1",
					CompanyName: req.CompanyName,
					Email:       req.Email,
					Status:      "pending",
				},
				AccessToken: "access-token-123",
			}, nil
		},
	}
	client, cleanup := setupTestServer(t, custSrv, nil)
	defer cleanup()

	adapter := NewCustomerServiceAdapter(client)
	resp, err := adapter.RegisterCustomer(context.Background(), &customerv1.RegisterCustomerRequest{
		Token:       "invite-token",
		CompanyName: "ACME Corp",
		Email:       "admin@acme.com",
		Password:    "secure-pass",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "new-cust-1", resp.GetCustomer().GetId())
	assert.Equal(t, "access-token-123", resp.GetAccessToken())
	assert.Equal(t, "pending", resp.GetCustomer().GetStatus())
}

func TestCustomerServiceAdapterRegisterCustomerError(t *testing.T) {
	custSrv := &testCustomerServer{
		registerCustomerFn: func(ctx context.Context, req *customerv1.RegisterCustomerRequest) (*customerv1.RegisterCustomerResponse, error) {
			return nil, status.Error(codes.InvalidArgument, "invalid token")
		},
	}
	client, cleanup := setupTestServer(t, custSrv, nil)
	defer cleanup()

	adapter := NewCustomerServiceAdapter(client)
	resp, err := adapter.RegisterCustomer(context.Background(), &customerv1.RegisterCustomerRequest{
		Token: "bad-token",
	})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, ErrInvalidArgument)
}

func TestCustomerServiceAdapterGetCustomer(t *testing.T) {
	custSrv := &testCustomerServer{
		getCustomerFn: func(ctx context.Context, req *customerv1.GetCustomerRequest) (*customerv1.GetCustomerResponse, error) {
			return &customerv1.GetCustomerResponse{
				Customer: &customerv1.Customer{
					Id:          req.Id,
					CompanyName: "ACME Corp",
					Email:       "admin@acme.com",
					Status:      "active",
				},
			}, nil
		},
	}
	client, cleanup := setupTestServer(t, custSrv, nil)
	defer cleanup()

	adapter := NewCustomerServiceAdapter(client)
	resp, err := adapter.GetCustomer(context.Background(), "cust-42")
	require.NoError(t, err)
	assert.Equal(t, "cust-42", resp.GetCustomer().GetId())
	assert.Equal(t, "ACME Corp", resp.GetCustomer().GetCompanyName())
}

func TestCustomerServiceAdapterGetCustomerNotFound(t *testing.T) {
	custSrv := &testCustomerServer{
		getCustomerFn: func(ctx context.Context, req *customerv1.GetCustomerRequest) (*customerv1.GetCustomerResponse, error) {
			return nil, status.Error(codes.NotFound, "customer not found")
		},
	}
	client, cleanup := setupTestServer(t, custSrv, nil)
	defer cleanup()

	adapter := NewCustomerServiceAdapter(client)
	resp, err := adapter.GetCustomer(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestCustomerServiceAdapterListCustomers(t *testing.T) {
	custSrv := &testCustomerServer{
		listCustomersFn: func(ctx context.Context, req *customerv1.ListCustomersRequest) (*customerv1.ListCustomersResponse, error) {
			return &customerv1.ListCustomersResponse{
				Customers: []*customerv1.Customer{
					{Id: "c1", CompanyName: "Co 1", Status: "active"},
					{Id: "c2", CompanyName: "Co 2", Status: "pending"},
				},
				Pagination: &commonv1.Pagination{
					Page: 1, PageSize: 10, Total: 2,
				},
			}, nil
		},
	}
	client, cleanup := setupTestServer(t, custSrv, nil)
	defer cleanup()

	adapter := NewCustomerServiceAdapter(client)
	resp, err := adapter.ListCustomers(context.Background(), &customerv1.ListCustomersRequest{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Len(t, resp.GetCustomers(), 2)
	assert.Equal(t, int32(1), resp.GetPagination().GetPage())
	assert.Equal(t, int32(2), resp.GetPagination().GetTotal())
}

func TestCustomerServiceAdapterUpdateCustomer(t *testing.T) {
	custSrv := &testCustomerServer{
		updateCustomerFn: func(ctx context.Context, req *customerv1.UpdateCustomerRequest) (*customerv1.UpdateCustomerResponse, error) {
			return &customerv1.UpdateCustomerResponse{
				Customer: &customerv1.Customer{
					Id:     req.Id,
					Phone:  req.Phone,
					City:   req.City,
					Status: "active",
				},
			}, nil
		},
	}
	client, cleanup := setupTestServer(t, custSrv, nil)
	defer cleanup()

	adapter := NewCustomerServiceAdapter(client)
	resp, err := adapter.UpdateCustomer(context.Background(), &customerv1.UpdateCustomerRequest{
		Id:    "cust-1",
		Phone: "555-0100",
		City:  "New York",
	})
	require.NoError(t, err)
	assert.Equal(t, "cust-1", resp.GetCustomer().GetId())
	assert.Equal(t, "555-0100", resp.GetCustomer().GetPhone())
}

func TestCustomerServiceAdapterApproveCustomer(t *testing.T) {
	custSrv := &testCustomerServer{
		approveCustomerFn: func(ctx context.Context, req *customerv1.ApproveCustomerRequest) (*customerv1.ApproveCustomerResponse, error) {
			return &customerv1.ApproveCustomerResponse{
				Customer: &customerv1.Customer{
					Id:     req.Id,
					Status: "active",
				},
			}, nil
		},
	}
	client, cleanup := setupTestServer(t, custSrv, nil)
	defer cleanup()

	adapter := NewCustomerServiceAdapter(client)
	resp, err := adapter.ApproveCustomer(context.Background(), "cust-1", "admin-user")
	require.NoError(t, err)
	assert.Equal(t, "cust-1", resp.GetCustomer().GetId())
	assert.Equal(t, "active", resp.GetCustomer().GetStatus())
}

func TestCustomerServiceAdapterRejectCustomer(t *testing.T) {
	custSrv := &testCustomerServer{
		rejectCustomerFn: func(ctx context.Context, req *customerv1.RejectCustomerRequest) (*customerv1.RejectCustomerResponse, error) {
			return &customerv1.RejectCustomerResponse{}, nil
		},
	}
	client, cleanup := setupTestServer(t, custSrv, nil)
	defer cleanup()

	adapter := NewCustomerServiceAdapter(client)
	_, err := adapter.RejectCustomer(context.Background(), "cust-1", "invalid docs")
	assert.NoError(t, err)
}

func TestCustomerServiceAdapterRejectCustomerError(t *testing.T) {
	custSrv := &testCustomerServer{
		rejectCustomerFn: func(ctx context.Context, req *customerv1.RejectCustomerRequest) (*customerv1.RejectCustomerResponse, error) {
			return nil, status.Error(codes.PermissionDenied, "no permission")
		},
	}
	client, cleanup := setupTestServer(t, custSrv, nil)
	defer cleanup()

	adapter := NewCustomerServiceAdapter(client)
	_, err := adapter.RejectCustomer(context.Background(), "cust-1", "reason")
	assert.ErrorIs(t, err, ErrPermissionDenied)
}

func TestCustomerServiceAdapterSuspendCustomer(t *testing.T) {
	custSrv := &testCustomerServer{
		suspendCustomerFn: func(ctx context.Context, req *customerv1.SuspendCustomerRequest) (*customerv1.SuspendCustomerResponse, error) {
			return &customerv1.SuspendCustomerResponse{}, nil
		},
	}
	client, cleanup := setupTestServer(t, custSrv, nil)
	defer cleanup()

	adapter := NewCustomerServiceAdapter(client)
	_, err := adapter.SuspendCustomer(context.Background(), "cust-1", "violation")
	assert.NoError(t, err)
}

func TestCustomerServiceAdapterAssignCompanyRole(t *testing.T) {
	permSrv := &testPermissionServer{
		assignCompanyRoleFn: func(ctx context.Context, req *customerv1.AssignCompanyRoleRequest) (*customerv1.AssignCompanyRoleResponse, error) {
			return &customerv1.AssignCompanyRoleResponse{}, nil
		},
	}
	client, cleanup := setupTestServer(t, nil, permSrv)
	defer cleanup()

	adapter := NewCustomerServiceAdapter(client)
	_, err := adapter.AssignCompanyRole(context.Background(), &customerv1.AssignCompanyRoleRequest{
		CustomerId: "cust-1",
		UserId:     "user-1",
		RoleName:   "admin",
	})
	assert.NoError(t, err)
}

func TestCustomerServiceAdapterListCompanyRoles(t *testing.T) {
	permSrv := &testPermissionServer{
		listCompanyRolesFn: func(ctx context.Context, req *customerv1.ListCompanyRolesRequest) (*customerv1.ListCompanyRolesResponse, error) {
			return &customerv1.ListCompanyRolesResponse{
				Roles: []*customerv1.CompanyRole{
					{Id: "r1", Name: "admin", CustomerId: req.CustomerId, Permissions: []string{"read", "write"}, IsDefault: false},
				},
			}, nil
		},
	}
	client, cleanup := setupTestServer(t, nil, permSrv)
	defer cleanup()

	adapter := NewCustomerServiceAdapter(client)
	resp, err := adapter.ListCompanyRoles(context.Background(), &customerv1.ListCompanyRolesRequest{
		CustomerId: "cust-1",
	})
	require.NoError(t, err)
	assert.Len(t, resp.GetRoles(), 1)
	assert.Equal(t, "admin", resp.GetRoles()[0].GetName())
}

func TestCustomerServiceAdapterGetUserPermissions(t *testing.T) {
	permSrv := &testPermissionServer{
		getUserPermissionsFn: func(ctx context.Context, req *customerv1.GetUserPermissionsRequest) (*customerv1.GetUserPermissionsResponse, error) {
			return &customerv1.GetUserPermissionsResponse{
				Permissions: []string{"customer:read", "customer:write"},
			}, nil
		},
	}
	client, cleanup := setupTestServer(t, nil, permSrv)
	defer cleanup()

	adapter := NewCustomerServiceAdapter(client)
	resp, err := adapter.GetUserPermissions(context.Background(), &customerv1.GetUserPermissionsRequest{
		UserId:     "user-1",
		CustomerId: "cust-1",
	})
	require.NoError(t, err)
	assert.Len(t, resp.GetPermissions(), 2)
	assert.Equal(t, "customer:read", resp.GetPermissions()[0])
}

func TestGRPCDialErrorPassthrough(t *testing.T) {
	client, err := NewCustomerClient("127.0.0.1:1")
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	adapter := NewCustomerServiceAdapter(client)
	resp, err := adapter.GetCustomer(context.Background(), "test")
	assert.Error(t, err)
	assert.Nil(t, resp)
}
