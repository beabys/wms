package grpc_test

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"

	"github.com/beabys/wms/customer-service/internal/application/customer/usecase"
	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
	grpcadapter "github.com/beabys/wms/customer-service/internal/infrastructure/adapters/grpc"
)

// ---------------------------------------------------------------------------
// Mock repositories
// ---------------------------------------------------------------------------

type mockCustomerRepo struct {
	getByIDFn func(ctx context.Context, id string) (*model.Customer, error)
	saveFn    func(ctx context.Context, c *model.Customer) error
	updateFn  func(ctx context.Context, c *model.Customer) error
	listFn    func(ctx context.Context, status string, page, pageSize int32) ([]*model.Customer, int32, error)
}

func (m *mockCustomerRepo) GetByID(ctx context.Context, id string) (*model.Customer, error) { return m.getByIDFn(ctx, id) }
func (m *mockCustomerRepo) GetByEmail(ctx context.Context, email string) (*model.Customer, error) { return nil, nil }
func (m *mockCustomerRepo) Save(ctx context.Context, c *model.Customer) error             { return m.saveFn(ctx, c) }
func (m *mockCustomerRepo) Update(ctx context.Context, c *model.Customer) error           { return m.updateFn(ctx, c) }
func (m *mockCustomerRepo) List(ctx context.Context, status string, page, pageSize int32) ([]*model.Customer, int32, error) {
	return m.listFn(ctx, status, page, pageSize)
}

type mockInviteRepo struct {
	createFn func(ctx context.Context, link *model.InviteLink) error
}

func (m *mockInviteRepo) Create(ctx context.Context, link *model.InviteLink) error { return m.createFn(ctx, link) }
func (m *mockInviteRepo) GetByToken(ctx context.Context, token string) (*model.InviteLink, error) { return nil, nil }
func (m *mockInviteRepo) MarkUsed(ctx context.Context, id string) error { return nil }

// ---------------------------------------------------------------------------
// Test setup factory
// ---------------------------------------------------------------------------

func startGRPCServer(t *testing.T, customerMock *mockCustomerRepo, inviteMock *mockInviteRepo) (customerv1.CustomerServiceClient, func()) {
	t.Helper()

	log := zap.NewNop()
	listener := bufconn.Listen(1024 * 1024)

	customerService := usecase.NewCustomerService(customerMock, inviteMock)
	srv := grpcadapter.NewServer(log, customerService)

	gsrv := grpc.NewServer()
	customerv1.RegisterCustomerServiceServer(gsrv, srv)
	go gsrv.Serve(listener)

	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	return customerv1.NewCustomerServiceClient(conn), func() { conn.Close(); gsrv.Stop(); listener.Close() }
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestGRPCCreateCustomer(t *testing.T) {
	t.Parallel()

	addr, _ := model.NewAddress("Str 1", "", "Berlin", "10115", "DE")
	customerRepo := &mockCustomerRepo{
		saveFn: func(ctx context.Context, c *model.Customer) error {
			c.ID = "new-id"
			c.Status = model.CustomerStatusPending
			c.Address = addr
			return nil
		},
	}

	client, cleanup := startGRPCServer(t, customerRepo, &mockInviteRepo{
		createFn: func(ctx context.Context, link *model.InviteLink) error {
			link.ID = "inv-1"
			link.Token = "tok_abc123"
			return nil
		},
	})
	defer cleanup()

	resp, err := client.CreateCustomer(context.Background(), &customerv1.CreateCustomerRequest{
		CompanyName: "New GmbH",
		VatNumber:   "DE123456789",
		Address: &commonv1.Address{
			Line1:      "Str 1",
			City:       "Berlin",
			PostalCode: "10115",
			Country:    "DE",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.GetCustomer())
	assert.Equal(t, "New GmbH", resp.GetCustomer().GetCompanyName())
	assert.Equal(t, "pending", resp.GetCustomer().GetStatus())
}

func TestGRPCGetCustomer(t *testing.T) {
	t.Parallel()

	addr, _ := model.NewAddress("Str 1", "", "Berlin", "10115", "DE")
	customerRepo := &mockCustomerRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.Customer, error) {
			return &model.Customer{ID: id, CompanyName: "Known GmbH", Status: model.CustomerStatusActive, Address: addr}, nil
		},
	}

	client, cleanup := startGRPCServer(t, customerRepo, &mockInviteRepo{createFn: func(ctx context.Context, link *model.InviteLink) error { return nil }})
	defer cleanup()

	resp, err := client.GetCustomer(context.Background(), &customerv1.GetCustomerRequest{Id: "cust_1"})
	require.NoError(t, err)
	assert.Equal(t, "cust_1", resp.GetCustomer().GetId())
	assert.Equal(t, "Known GmbH", resp.GetCustomer().GetCompanyName())
}

func TestGRPCListCustomers(t *testing.T) {
	t.Parallel()

	addr, _ := model.NewAddress("Str 1", "", "Berlin", "10115", "DE")
	customerRepo := &mockCustomerRepo{
		listFn: func(ctx context.Context, status string, page, pageSize int32) ([]*model.Customer, int32, error) {
			return []*model.Customer{
				{ID: "1", CompanyName: "Alpha", Address: addr, Status: model.CustomerStatusActive},
			}, 1, nil
		},
	}

	client, cleanup := startGRPCServer(t, customerRepo, &mockInviteRepo{createFn: func(ctx context.Context, link *model.InviteLink) error { return nil }})
	defer cleanup()

	resp, err := client.ListCustomers(context.Background(), &customerv1.ListCustomersRequest{
		Pagination: &commonv1.Pagination{Page: 1, Limit: 10},
	})
	require.NoError(t, err)
	assert.Len(t, resp.GetCustomers(), 1)
	assert.Equal(t, int32(1), resp.GetPagination().GetTotal())
}

func TestGRPCApproveCustomer(t *testing.T) {
	t.Parallel()

	addr, _ := model.NewAddress("Str 1", "", "Berlin", "10115", "DE")
	customerRepo := &mockCustomerRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.Customer, error) {
			return &model.Customer{ID: id, CompanyName: "Pending GmbH", Status: model.CustomerStatusPending, Address: addr}, nil
		},
		updateFn: func(ctx context.Context, c *model.Customer) error { return nil },
	}

	client, cleanup := startGRPCServer(t, customerRepo, &mockInviteRepo{createFn: func(ctx context.Context, link *model.InviteLink) error { return nil }})
	defer cleanup()

	resp, err := client.ApproveCustomer(context.Background(), &customerv1.ApproveCustomerRequest{CustomerId: "cust_1"})
	require.NoError(t, err)
	assert.Equal(t, "active", resp.GetCustomer().GetStatus())
}

func TestGRPCSuspendCustomer(t *testing.T) {
	t.Parallel()

	addr, _ := model.NewAddress("Str 1", "", "Berlin", "10115", "DE")
	customerRepo := &mockCustomerRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.Customer, error) {
			return &model.Customer{ID: id, CompanyName: "Active GmbH", Status: model.CustomerStatusActive, Address: addr}, nil
		},
		updateFn: func(ctx context.Context, c *model.Customer) error { return nil },
	}

	client, cleanup := startGRPCServer(t, customerRepo, &mockInviteRepo{createFn: func(ctx context.Context, link *model.InviteLink) error { return nil }})
	defer cleanup()

	resp, err := client.SuspendCustomer(context.Background(), &customerv1.SuspendCustomerRequest{
		CustomerId: "cust_1",
		Reason:     "non-payment",
	})
	require.NoError(t, err)
	assert.Equal(t, "suspended", resp.GetCustomer().GetStatus())
}

func TestGRPCInviteCustomer(t *testing.T) {
	t.Parallel()

	customerRepo := &mockCustomerRepo{}
	inviteRepo := &mockInviteRepo{
		createFn: func(ctx context.Context, link *model.InviteLink) error {
			link.ID = "inv-1"
			link.Token = "tok_abc123"
			return nil
		},
	}

	client, cleanup := startGRPCServer(t, customerRepo, inviteRepo)
	defer cleanup()

	resp, err := client.InviteCustomer(context.Background(), &customerv1.InviteCustomerRequest{Email: "test@example.com"})
	require.NoError(t, err)
	assert.Equal(t, "tok_abc123", resp.GetInvitationLink())
}
