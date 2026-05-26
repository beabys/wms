package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/beabys/wms/customer-service/internal/application/customer/command"
	"github.com/beabys/wms/customer-service/internal/application/customer/usecase"
	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
)

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

type mockCustomerRepo struct {
	getByIDFn  func(ctx context.Context, id string) (*model.Customer, error)
	getByEmail func(ctx context.Context, email string) (*model.Customer, error)
	saveFn     func(ctx context.Context, c *model.Customer) error
	listFn     func(ctx context.Context, status string, page, pageSize int32) ([]*model.Customer, int32, error)
	updateFn   func(ctx context.Context, c *model.Customer) error
}

func (m *mockCustomerRepo) GetByID(ctx context.Context, id string) (*model.Customer, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockCustomerRepo) GetByEmail(ctx context.Context, email string) (*model.Customer, error) {
	return m.getByEmail(ctx, email)
}
func (m *mockCustomerRepo) Save(ctx context.Context, c *model.Customer) error {
	return m.saveFn(ctx, c)
}
func (m *mockCustomerRepo) List(ctx context.Context, status string, page, pageSize int32) ([]*model.Customer, int32, error) {
	return m.listFn(ctx, status, page, pageSize)
}
func (m *mockCustomerRepo) Update(ctx context.Context, c *model.Customer) error {
	return m.updateFn(ctx, c)
}

type mockInviteRepo struct {
	createFn   func(ctx context.Context, link *model.InviteLink) error
	getByToken func(ctx context.Context, token string) (*model.InviteLink, error)
	markUsedFn func(ctx context.Context, id string) error
}

func (m *mockInviteRepo) Create(ctx context.Context, link *model.InviteLink) error {
	return m.createFn(ctx, link)
}
func (m *mockInviteRepo) GetByToken(ctx context.Context, token string) (*model.InviteLink, error) {
	return m.getByToken(ctx, token)
}
func (m *mockInviteRepo) MarkUsed(ctx context.Context, id string) error {
	return m.markUsedFn(ctx, id)
}

// noopInviteRepo returns a minimal invite repo that succeeds for Create and returns nil for others.
func noopInviteRepo() *mockInviteRepo {
	return &mockInviteRepo{
		createFn: func(ctx context.Context, link *model.InviteLink) error { return nil },
		getByToken: func(ctx context.Context, token string) (*model.InviteLink, error) { return nil, nil },
		markUsedFn: func(ctx context.Context, id string) error { return nil },
	}
}

// noopCustomerRepo returns a minimal customer repo that returns nil for everything.
func noopCustomerRepo() *mockCustomerRepo {
	return &mockCustomerRepo{
		getByIDFn:  func(ctx context.Context, id string) (*model.Customer, error) { return nil, errors.New("not found") },
		getByEmail: func(ctx context.Context, email string) (*model.Customer, error) { return nil, nil },
		saveFn:     func(ctx context.Context, c *model.Customer) error { return nil },
		listFn:     func(ctx context.Context, status string, page, pageSize int32) ([]*model.Customer, int32, error) { return nil, 0, nil },
		updateFn:   func(ctx context.Context, c *model.Customer) error { return nil },
	}
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestCreateCustomerHandler_Success(t *testing.T) {
	t.Parallel()

	var saved *model.Customer
	repo := &mockCustomerRepo{
		saveFn: func(ctx context.Context, c *model.Customer) error {
			saved = c
			return nil
		},
	}

	svc := usecase.NewCustomerService(repo, noopInviteRepo())
	addr, _ := model.NewAddress("Strasse 1", "", "Berlin", "10115", "DE")
	cmd := command.CreateCustomerCommand{
		CompanyName: "Test GmbH",
		VATNumber:   "DE123456789",
		Address:     addr,
		RateCardID:  "rate_1",
	}

	customer, err := svc.CreateCustomer(context.Background(), &cmd)
	require.NoError(t, err)
	require.NotNil(t, customer)

	assert.Equal(t, "Test GmbH", customer.CompanyName)
	assert.Equal(t, "DE123456789", customer.VATNumber)
	assert.Equal(t, model.CustomerStatusPending, customer.Status)
	assert.NotEmpty(t, customer.ID)
	assert.NotEmpty(t, saved)
}

func TestCreateCustomerHandler_InvalidVAT(t *testing.T) {
	t.Parallel()

	repo := noopCustomerRepo()
	svc := usecase.NewCustomerService(repo, noopInviteRepo())

	addr, _ := model.NewAddress("Strasse 1", "", "Berlin", "10115", "DE")
	cmd := command.CreateCustomerCommand{
		CompanyName: "Bad VAT GmbH",
		VATNumber:   "INVALID",
		Address:     addr,
	}

	_, err := svc.CreateCustomer(context.Background(), &cmd)
	assert.Error(t, err)
}

func TestInviteCustomerHandler_Success(t *testing.T) {
	t.Parallel()

	var saved *model.InviteLink
	repo := &mockInviteRepo{
		createFn: func(ctx context.Context, link *model.InviteLink) error {
			saved = link
			return nil
		},
	}

	svc := usecase.NewCustomerService(noopCustomerRepo(), repo)
	cmd := command.InviteCustomerCommand{Email: "test@example.com"}

	invite, err := svc.InviteCustomer(context.Background(), &cmd)
	require.NoError(t, err)
	require.NotNil(t, invite)

	assert.Equal(t, "test@example.com", invite.Email)
	assert.False(t, invite.Used)
	assert.True(t, invite.ExpiresAt.After(time.Now()))
	assert.NotEmpty(t, invite.Token)
	assert.NotEmpty(t, saved)
}

func TestInviteCustomerHandler_EmptyEmail(t *testing.T) {
	t.Parallel()

	repo := noopInviteRepo()
	svc := usecase.NewCustomerService(noopCustomerRepo(), repo)

	_, err := svc.InviteCustomer(context.Background(), &command.InviteCustomerCommand{Email: ""})
	assert.Error(t, err)
}

func TestApproveCustomerHandler_Success(t *testing.T) {
	t.Parallel()

	repo := &mockCustomerRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.Customer, error) {
			return &model.Customer{
				ID:     id,
				Status: model.CustomerStatusPending,
			}, nil
		},
		updateFn: func(ctx context.Context, c *model.Customer) error {
			return nil
		},
	}

	svc := usecase.NewCustomerService(repo, noopInviteRepo())
	cmd := command.ApproveCustomerCommand{CustomerID: "cust_1"}

	customer, err := svc.ApproveCustomer(context.Background(), &cmd)
	require.NoError(t, err)
	assert.Equal(t, model.CustomerStatusActive, customer.Status)
}

func TestApproveCustomerHandler_AlreadyActive(t *testing.T) {
	t.Parallel()

	repo := &mockCustomerRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.Customer, error) {
			return &model.Customer{
				ID:     id,
				Status: model.CustomerStatusActive,
			}, nil
		},
	}

	svc := usecase.NewCustomerService(repo, noopInviteRepo())
	_, err := svc.ApproveCustomer(context.Background(), &command.ApproveCustomerCommand{CustomerID: "cust_1"})
	assert.Error(t, err)
}

func TestSuspendCustomerHandler_Success(t *testing.T) {
	t.Parallel()

	repo := &mockCustomerRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.Customer, error) {
			return &model.Customer{
				ID:     id,
				Status: model.CustomerStatusActive,
			}, nil
		},
		updateFn: func(ctx context.Context, c *model.Customer) error {
			return nil
		},
	}

	svc := usecase.NewCustomerService(repo, noopInviteRepo())
	cmd := command.SuspendCustomerCommand{CustomerID: "cust_1", Reason: "non-payment"}

	customer, err := svc.SuspendCustomer(context.Background(), &cmd)
	require.NoError(t, err)
	assert.Equal(t, model.CustomerStatusSuspended, customer.Status)
}

func TestSuspendCustomerHandler_NotActive(t *testing.T) {
	t.Parallel()

	repo := &mockCustomerRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.Customer, error) {
			return &model.Customer{
				ID:     id,
				Status: model.CustomerStatusPending,
			}, nil
		},
	}

	svc := usecase.NewCustomerService(repo, noopInviteRepo())
	_, err := svc.SuspendCustomer(context.Background(), &command.SuspendCustomerCommand{CustomerID: "cust_1"})
	assert.Error(t, err)
}

func TestGetCustomerHandler_Success(t *testing.T) {
	t.Parallel()

	repo := &mockCustomerRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.Customer, error) {
			return &model.Customer{ID: id, CompanyName: "Known GmbH"}, nil
		},
	}

	svc := usecase.NewCustomerService(repo, noopInviteRepo())
	customer, err := svc.GetCustomer(context.Background(), &command.GetCustomerQuery{CustomerID: "cust_1"})
	require.NoError(t, err)
	assert.Equal(t, "Known GmbH", customer.CompanyName)
}

func TestGetCustomerHandler_NotFound(t *testing.T) {
	t.Parallel()

	repo := &mockCustomerRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.Customer, error) {
			return nil, errors.New("customer not found")
		},
	}

	svc := usecase.NewCustomerService(repo, noopInviteRepo())
	_, err := svc.GetCustomer(context.Background(), &command.GetCustomerQuery{CustomerID: "missing"})
	assert.Error(t, err)
}

func TestListCustomersHandler_Success(t *testing.T) {
	t.Parallel()

	repo := &mockCustomerRepo{
		listFn: func(ctx context.Context, status string, page, pageSize int32) ([]*model.Customer, int32, error) {
			return []*model.Customer{
				{ID: "1", CompanyName: "Alpha"},
				{ID: "2", CompanyName: "Beta"},
			}, 2, nil
		},
	}

	svc := usecase.NewCustomerService(repo, noopInviteRepo())
	customers, total, err := svc.ListCustomers(context.Background(), &command.ListCustomersQuery{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Len(t, customers, 2)
	assert.Equal(t, int32(2), total)
}
