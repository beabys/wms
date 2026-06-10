package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/beabys/wms/customer-service/internal/application/customer/command"
	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
	repomocks "github.com/beabys/wms/customer-service/mocks/application/customer/repository"
	portmocks "github.com/beabys/wms/customer-service/mocks/application/customer/ports"
	"github.com/beabys/wms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Tests ---

func TestRegisterCustomer_Success(t *testing.T) {
	customerRepo := repomocks.NewCustomerRepository(t)
	companyRepo := repomocks.NewCompanyRepository(t)
	roleRepo := repomocks.NewCompanyRoleRepository(t)
	auditLogRepo := repomocks.NewAuditLogRepository(t)
	authClient := portmocks.NewAuthClient(t)

	authClient.On("ValidateInvite", mock.Anything, "valid-token").Return("admin@testcorp.com", "admin-123", nil)
	authClient.On("CreateUser", mock.Anything, mock.Anything).Return(&command.CreateUserResponse{
		UserID: "auth-user-123",
		Email:  "admin@testcorp.com",
		Name:   "Test Corp Admin",
		Role:   "customer_admin",
	}, nil)
	customerRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Customer")).Return(nil)
	companyRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Company")).Return(nil)
	roleRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.CompanyRole")).Return(nil)
	roleRepo.On("AssignUserRole", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	uc := NewCustomerUseCase(&testLogger{}, customerRepo, companyRepo, roleRepo, auditLogRepo, authClient)

	result, err := uc.RegisterCustomer(context.Background(), command.RegisterCustomerCommand{
		Token:       "valid-token",
		CompanyName: "Test Corp",
		Email:       "admin@testcorp.com",
		Password:    "password123",
		Phone:       "1234567890",
		VatNumber:   "VAT123",
		Address:     "123 Main St",
		City:        "New York",
		PostalCode:  "10001",
		Country:     "US",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.CustomerID)
}

func TestGetCustomer_Success(t *testing.T) {
	customerRepo := repomocks.NewCustomerRepository(t)
	now := time.Now()

	customerRepo.On("GetByID", mock.Anything, "cust-1").Return(&model.Customer{
		ID:          "cust-1",
		CompanyName: "Test Corp",
		Email:       "admin@testcorp.com",
		Status:      model.CustomerStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil)

	uc := NewCustomerUseCase(&testLogger{}, customerRepo, nil, nil, nil, nil)

	result, err := uc.GetCustomer(context.Background(), "cust-1")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "cust-1", result.ID)
	assert.Equal(t, "Test Corp", result.CompanyName)
	assert.Equal(t, "active", result.Status)
}

func TestGetCustomer_NotFound(t *testing.T) {
	customerRepo := repomocks.NewCustomerRepository(t)

	customerRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, nil)

	uc := NewCustomerUseCase(&testLogger{}, customerRepo, nil, nil, nil, nil)

	_, err := uc.GetCustomer(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "customer not found")
}

func TestGetCustomer_EmptyID(t *testing.T) {
	uc := &CustomerUseCase{logger: &testLogger{}}
	_, err := uc.GetCustomer(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "customer ID is required")
}

func TestApproveCustomer_Success(t *testing.T) {
	customerRepo := repomocks.NewCustomerRepository(t)
	auditLogRepo := repomocks.NewAuditLogRepository(t)
	now := time.Now()

	customerRepo.On("GetByID", mock.Anything, "cust-1").Return(&model.Customer{
		ID:        "cust-1",
		Status:    model.CustomerStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)
	customerRepo.On("UpdateStatus", mock.Anything, "cust-1", model.CustomerStatusActive, "").Return(nil)
	auditLogRepo.On("Insert", mock.Anything, mock.AnythingOfType("*model.AuditLog")).Return(nil)

	uc := NewCustomerUseCase(&testLogger{}, customerRepo, nil, nil, auditLogRepo, nil)

	result, err := uc.ApproveCustomer(context.Background(), command.ApproveCustomerCommand{
		CustomerID: "cust-1",
		ApprovedBy: "admin-123",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "active", result.Status)
}

func TestApproveCustomer_NotPending(t *testing.T) {
	customerRepo := repomocks.NewCustomerRepository(t)
	now := time.Now()

	customerRepo.On("GetByID", mock.Anything, "cust-1").Return(&model.Customer{
		ID:        "cust-1",
		Status:    model.CustomerStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	uc := NewCustomerUseCase(&testLogger{}, customerRepo, nil, nil, nil, nil)

	_, err := uc.ApproveCustomer(context.Background(), command.ApproveCustomerCommand{
		CustomerID: "cust-1",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not in pending status")
}

func TestRejectCustomer_Success(t *testing.T) {
	customerRepo := repomocks.NewCustomerRepository(t)
	auditLogRepo := repomocks.NewAuditLogRepository(t)
	now := time.Now()

	customerRepo.On("GetByID", mock.Anything, "cust-1").Return(&model.Customer{
		ID:        "cust-1",
		Status:    model.CustomerStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)
	customerRepo.On("UpdateStatus", mock.Anything, "cust-1", model.CustomerStatusRejected, "Invalid documentation").Return(nil)
	auditLogRepo.On("Insert", mock.Anything, mock.AnythingOfType("*model.AuditLog")).Return(nil)

	uc := NewCustomerUseCase(&testLogger{}, customerRepo, nil, nil, auditLogRepo, nil)

	result, err := uc.RejectCustomer(context.Background(), command.RejectCustomerCommand{
		CustomerID: "cust-1",
		Reason:     "Invalid documentation",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "rejected", result.Status)
}

func TestRejectCustomer_EmptyReason(t *testing.T) {
	uc := &CustomerUseCase{logger: &testLogger{}}
	_, err := uc.RejectCustomer(context.Background(), command.RejectCustomerCommand{
		CustomerID: "cust-1",
		Reason:     "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reject reason is required")
}

func TestSuspendCustomer_Success(t *testing.T) {
	customerRepo := repomocks.NewCustomerRepository(t)
	auditLogRepo := repomocks.NewAuditLogRepository(t)
	now := time.Now()

	customerRepo.On("GetByID", mock.Anything, "cust-1").Return(&model.Customer{
		ID:        "cust-1",
		Status:    model.CustomerStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)
	customerRepo.On("UpdateStatus", mock.Anything, "cust-1", model.CustomerStatusSuspended, "Policy violation").Return(nil)
	auditLogRepo.On("Insert", mock.Anything, mock.AnythingOfType("*model.AuditLog")).Return(nil)

	uc := NewCustomerUseCase(&testLogger{}, customerRepo, nil, nil, auditLogRepo, nil)

	result, err := uc.SuspendCustomer(context.Background(), command.SuspendCustomerCommand{
		CustomerID: "cust-1",
		Reason:     "Policy violation",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "suspended", result.Status)
}

func TestRestoreCustomer_Success(t *testing.T) {
	customerRepo := repomocks.NewCustomerRepository(t)
	auditLogRepo := repomocks.NewAuditLogRepository(t)
	now := time.Now()

	customerRepo.On("GetByID", mock.Anything, "cust-1").Return(&model.Customer{
		ID:        "cust-1",
		Status:    model.CustomerStatusSuspended,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)
	customerRepo.On("UpdateStatus", mock.Anything, "cust-1", model.CustomerStatusActive, "").Return(nil)
	auditLogRepo.On("Insert", mock.Anything, mock.AnythingOfType("*model.AuditLog")).Return(nil)

	uc := NewCustomerUseCase(&testLogger{}, customerRepo, nil, nil, auditLogRepo, nil)

	result, err := uc.RestoreCustomer(context.Background(), command.RestoreCustomerCommand{
		CustomerID: "cust-1",
		RestoredBy: "admin-123",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "active", result.Status)
}

func TestRestoreCustomer_NotSuspended(t *testing.T) {
	customerRepo := repomocks.NewCustomerRepository(t)
	now := time.Now()

	customerRepo.On("GetByID", mock.Anything, "cust-1").Return(&model.Customer{
		ID:        "cust-1",
		Status:    model.CustomerStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	uc := NewCustomerUseCase(&testLogger{}, customerRepo, nil, nil, nil, nil)

	_, err := uc.RestoreCustomer(context.Background(), command.RestoreCustomerCommand{
		CustomerID: "cust-1",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not in suspended status")
}

func TestRestoreCustomer_NotFound(t *testing.T) {
	customerRepo := repomocks.NewCustomerRepository(t)

	customerRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, nil)

	uc := NewCustomerUseCase(&testLogger{}, customerRepo, nil, nil, nil, nil)

	_, err := uc.RestoreCustomer(context.Background(), command.RestoreCustomerCommand{
		CustomerID: "nonexistent",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "customer not found")
}

func TestRestoreCustomer_EmptyID(t *testing.T) {
	uc := &CustomerUseCase{logger: &testLogger{}}
	_, err := uc.RestoreCustomer(context.Background(), command.RestoreCustomerCommand{
		CustomerID: "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "customer ID is required")
}

func TestUpdateCustomer_Success(t *testing.T) {
	customerRepo := repomocks.NewCustomerRepository(t)
	now := time.Now()

	customerRepo.On("GetByID", mock.Anything, "cust-1").Return(&model.Customer{
		ID:          "cust-1",
		CompanyName: "Old Name",
		Email:       "admin@testcorp.com",
		Phone:       "",
		Status:      model.CustomerStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil)
	customerRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Customer")).Return(nil)

	uc := NewCustomerUseCase(&testLogger{}, customerRepo, nil, nil, nil, nil)

	result, err := uc.UpdateCustomer(context.Background(), command.UpdateCustomerCommand{
		CustomerID:  "cust-1",
		CompanyName: "New Name",
		Phone:       "9876543210",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "New Name", result.CompanyName)
	assert.Equal(t, "9876543210", result.Phone)
}

func TestListCustomers_Success(t *testing.T) {
	customerRepo := repomocks.NewCustomerRepository(t)
	now := time.Now()

	customerRepo.On("List", mock.Anything, 1, 20, "").Return([]*model.Customer{
		{ID: "cust-1", Status: model.CustomerStatusActive, CreatedAt: now, UpdatedAt: now},
		{ID: "cust-2", Status: model.CustomerStatusPending, CreatedAt: now, UpdatedAt: now},
	}, 2, nil)

	uc := NewCustomerUseCase(&testLogger{}, customerRepo, nil, nil, nil, nil)

	result, err := uc.ListCustomers(context.Background(), command.ListCustomersQuery{
		Page:     1,
		PageSize: 20,
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.TotalCount)
	assert.Len(t, result.Customers, 2)
}

func TestListCustomers_FilterByStatus(t *testing.T) {
	customerRepo := repomocks.NewCustomerRepository(t)
	now := time.Now()

	customerRepo.On("List", mock.Anything, 1, 20, "active").Return([]*model.Customer{
		{ID: "cust-1", Status: model.CustomerStatusActive, CreatedAt: now, UpdatedAt: now},
	}, 1, nil)

	uc := NewCustomerUseCase(&testLogger{}, customerRepo, nil, nil, nil, nil)

	result, err := uc.ListCustomers(context.Background(), command.ListCustomersQuery{
		Page:     1,
		PageSize: 20,
		Status:   "active",
	})
	require.NoError(t, err)
	assert.Equal(t, 1, result.TotalCount)
	assert.Len(t, result.Customers, 1)
	assert.Equal(t, "active", result.Customers[0].Status)
}

func TestListCustomers_Defaults(t *testing.T) {
	customerRepo := repomocks.NewCustomerRepository(t)

	customerRepo.On("List", mock.Anything, 1, 20, "").Return([]*model.Customer{}, 0, nil)

	uc := NewCustomerUseCase(&testLogger{}, customerRepo, nil, nil, nil, nil)

	result, err := uc.ListCustomers(context.Background(), command.ListCustomersQuery{})
	require.NoError(t, err)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 20, result.PageSize)
}

// --- Audit log tests ---

func TestGetAuditLogs_Success(t *testing.T) {
	auditLogRepo := repomocks.NewAuditLogRepository(t)

	auditLogRepo.On("ListByCustomerID", mock.Anything, "cust-1", 1, 20).Return([]*model.AuditLog{
		{
			ID:          "audit-1",
			CustomerID:  "cust-1",
			Action:      model.AuditActionApproved,
			PerformedBy: "admin-1",
			CreatedAt:   time.Now(),
		},
	}, 1, nil)

	uc := NewCustomerUseCase(&testLogger{}, nil, nil, nil, auditLogRepo, nil)

	result, err := uc.GetAuditLogs(context.Background(), command.ListAuditLogsQuery{
		CustomerID: "cust-1",
		Page:       1,
		PageSize:   20,
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, len(result.Entries))
	assert.Equal(t, "audit-1", result.Entries[0].ID)
	assert.Equal(t, "cust-1", result.Entries[0].CustomerID)
	assert.Equal(t, "approved", result.Entries[0].Action)
	assert.Equal(t, "admin-1", result.Entries[0].PerformedBy)
}

func TestGetAuditLogs_EmptyID(t *testing.T) {
	uc := &CustomerUseCase{logger: &testLogger{}}
	_, err := uc.GetAuditLogs(context.Background(), command.ListAuditLogsQuery{
		CustomerID: "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "customer ID is required")
}

func TestGetAuditLogs_Defaults(t *testing.T) {
	auditLogRepo := repomocks.NewAuditLogRepository(t)

	auditLogRepo.On("ListByCustomerID", mock.Anything, "cust-1", 1, 20).Return([]*model.AuditLog{}, 0, nil)

	uc := NewCustomerUseCase(&testLogger{}, nil, nil, nil, auditLogRepo, nil)

	result, err := uc.GetAuditLogs(context.Background(), command.ListAuditLogsQuery{
		CustomerID: "cust-1",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 20, result.PageSize)
}

// testLogger implements logger.Logger for use in tests.
type testLogger struct{}

func (l *testLogger) GetLogger() any                                    { return nil }
func (l *testLogger) Debug(msg string, fields ...logger.LogField)       {}
func (l *testLogger) Info(msg string, fields ...logger.LogField)        {}
func (l *testLogger) Warn(msg string, fields ...logger.LogField)        {}
func (l *testLogger) Error(msg string, err error, fields ...logger.LogField) {}
func (l *testLogger) Fatal(msg string, fields ...logger.LogField)       {}
