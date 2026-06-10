package usecase

import (
	"context"
	"testing"

	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"
	commonv1 "github.com/beabys/wms/proto/gen/go/common/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	mocks "github.com/beabys/wms/customer-service-bff/mocks/application/customer/usecase"
	"github.com/beabys/wms/customer-service-bff/internal/domain/model"
	"github.com/beabys/wms/pkg/logger"
)

func setupUseCase(t *testing.T, mockSvc *mocks.GrpcClient) *CustomerUseCase {
	t.Helper()
	log, _ := logger.NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	return New(log, mockSvc)
}

// --- RegisterCustomer ---

func TestRegisterCustomer_Success(t *testing.T) {
	mockSvc := mocks.NewGrpcClient(t)
	mockSvc.EXPECT().RegisterCustomer(mock.Anything, mock.MatchedBy(func(req *customerv1.RegisterCustomerRequest) bool {
		return req.CompanyName == "ACME Corp" && req.Email == "admin@acme.com"
	})).
		Return(&customerv1.RegisterCustomerResponse{
			Customer:    &customerv1.Customer{Id: "new-id", CompanyName: "ACME Corp", Status: "pending"},
			AccessToken: "token-123",
		}, nil)
	uc := setupUseCase(t, mockSvc)

	resp, err := uc.RegisterCustomer(context.Background(), &model.RegisterCustomerRequest{
		Token:       "invite-token",
		CompanyName: "ACME Corp",
		Email:       "admin@acme.com",
		Password:    "secure-pass",
	})
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "new-id", resp.Customer.ID)
	assert.Equal(t, "token-123", resp.AccessToken)
}

func TestRegisterCustomer_ValidationError(t *testing.T) {
	mockSvc := mocks.NewGrpcClient(t)
	uc := setupUseCase(t, mockSvc)

	_, err := uc.RegisterCustomer(context.Background(), &model.RegisterCustomerRequest{
		Token: "",
	})
	assert.Error(t, err)
}

// --- GetMyCustomer ---

func TestGetMyCustomer_Success(t *testing.T) {
	mockSvc := mocks.NewGrpcClient(t)
	mockSvc.EXPECT().GetCustomerByAdminID(mock.Anything, "admin-1").
		Return(&customerv1.GetCustomerByAdminIDResponse{
			Customer: &customerv1.Customer{Id: "cust-1", CompanyName: "Mock Co", Status: "active", CompanyAdminId: "admin-1"},
		}, nil)
	uc := setupUseCase(t, mockSvc)

	resp, err := uc.GetMyCustomer(context.Background(), "admin-1")
	require.NoError(t, err)
	assert.Equal(t, "cust-1", resp.ID)
	assert.Equal(t, "admin-1", resp.CompanyAdminID)
}

// --- GetCustomer ---

func TestGetCustomer_Success(t *testing.T) {
	mockSvc := mocks.NewGrpcClient(t)
	mockSvc.EXPECT().GetCustomer(mock.Anything, "cust-1").
		Return(&customerv1.GetCustomerResponse{
			Customer: &customerv1.Customer{Id: "cust-1", CompanyName: "Mock Co", Status: "active"},
		}, nil)
	uc := setupUseCase(t, mockSvc)

	resp, err := uc.GetCustomer(context.Background(), "cust-1")
	require.NoError(t, err)
	assert.Equal(t, "cust-1", resp.ID)
}

// --- ListCustomers ---

func TestListCustomers_Success(t *testing.T) {
	mockSvc := mocks.NewGrpcClient(t)
	mockSvc.EXPECT().ListCustomers(mock.Anything, mock.AnythingOfType("*customerv1.ListCustomersRequest")).
		Return(&customerv1.ListCustomersResponse{
			Customers:  []*customerv1.Customer{{Id: "c1", CompanyName: "Co 1"}},
			Pagination: &commonv1.Pagination{Page: 1, PageSize: 20, Total: 1},
		}, nil)
	uc := setupUseCase(t, mockSvc)

	page := 1
	pageSize := 10
	resp, err := uc.ListCustomers(context.Background(), nil, &page, &pageSize)
	require.NoError(t, err)
	assert.Len(t, resp.Customers, 1)
	assert.Equal(t, 1, resp.Pagination.Page)
}

// --- UpdateCustomer ---

func TestUpdateCustomer_Success(t *testing.T) {
	mockSvc := mocks.NewGrpcClient(t)
	mockSvc.EXPECT().UpdateCustomer(mock.Anything, mock.MatchedBy(func(req *customerv1.UpdateCustomerRequest) bool {
		return req.Id == "cust-1" && req.City == "New York"
	})).
		Return(&customerv1.UpdateCustomerResponse{
			Customer: &customerv1.Customer{Id: "cust-1", City: "New York", Status: "active"},
		}, nil)
	uc := setupUseCase(t, mockSvc)

	resp, err := uc.UpdateCustomer(context.Background(), "cust-1", &model.UpdateCustomerRequest{
		City: "New York",
	})
	require.NoError(t, err)
	assert.Equal(t, "cust-1", resp.ID)
	assert.Equal(t, "New York", resp.City)
}

// --- ApproveCustomer ---

func TestApproveCustomer_Success(t *testing.T) {
	mockSvc := mocks.NewGrpcClient(t)
	mockSvc.EXPECT().ApproveCustomer(mock.Anything, "cust-1", "admin-1").
		Return(&customerv1.ApproveCustomerResponse{
			Customer: &customerv1.Customer{Id: "cust-1", Status: "active"},
		}, nil)
	uc := setupUseCase(t, mockSvc)

	resp, err := uc.ApproveCustomer(context.Background(), "cust-1", "admin-1")
	require.NoError(t, err)
	assert.Equal(t, "cust-1", resp.ID)
	assert.Equal(t, "active", resp.Status)
}

// --- RejectCustomer ---

func TestRejectCustomer_Success(t *testing.T) {
	mockSvc := mocks.NewGrpcClient(t)
	mockSvc.EXPECT().RejectCustomer(mock.Anything, "cust-1", "invalid docs").
		Return(&customerv1.RejectCustomerResponse{}, nil)
	uc := setupUseCase(t, mockSvc)

	err := uc.RejectCustomer(context.Background(), "cust-1", "invalid docs")
	assert.NoError(t, err)
}

// --- SuspendCustomer ---

func TestSuspendCustomer_Success(t *testing.T) {
	mockSvc := mocks.NewGrpcClient(t)
	mockSvc.EXPECT().SuspendCustomer(mock.Anything, "cust-1", "violation").
		Return(&customerv1.SuspendCustomerResponse{}, nil)
	uc := setupUseCase(t, mockSvc)

	err := uc.SuspendCustomer(context.Background(), "cust-1", "violation")
	assert.NoError(t, err)
}

// --- RestoreCustomer ---

func TestRestoreCustomer_Success(t *testing.T) {
	mockSvc := mocks.NewGrpcClient(t)
	mockSvc.EXPECT().RestoreCustomer(mock.Anything, "cust-1", "admin-1").
		Return(&customerv1.RestoreCustomerResponse{
			Customer: &customerv1.Customer{Id: "cust-1", Status: "active"},
		}, nil)
	uc := setupUseCase(t, mockSvc)

	resp, err := uc.RestoreCustomer(context.Background(), "cust-1", "admin-1")
	require.NoError(t, err)
	assert.Equal(t, "cust-1", resp.ID)
}

// --- ListAuditLogs ---

func TestListAuditLogs_Success(t *testing.T) {
	mockSvc := mocks.NewGrpcClient(t)
	mockSvc.EXPECT().ListAuditLogs(mock.Anything, mock.AnythingOfType("*customerv1.ListAuditLogsRequest")).
		Return(&customerv1.ListAuditLogsResponse{
			Entries: []*customerv1.AuditEntry{
				{Id: "audit-1", CustomerId: "cust-1", Action: "approved", PerformedBy: "admin-1", CreatedAt: 1700000000},
			},
			Pagination: &commonv1.Pagination{Page: 1, PageSize: 20, Total: 1},
		}, nil)
	uc := setupUseCase(t, mockSvc)

	resp, err := uc.ListAuditLogs(context.Background(), "cust-1", 1, 20)
	require.NoError(t, err)
	assert.Len(t, resp.Entries, 1)
	assert.Equal(t, "approved", resp.Entries[0].Action)
}

// --- AssignCompanyRole ---

func TestAssignCompanyRole_Success(t *testing.T) {
	mockSvc := mocks.NewGrpcClient(t)
	mockSvc.EXPECT().AssignCompanyRole(mock.Anything, mock.AnythingOfType("*customerv1.AssignCompanyRoleRequest")).
		Return(&customerv1.AssignCompanyRoleResponse{}, nil)
	uc := setupUseCase(t, mockSvc)

	err := uc.AssignCompanyRole(context.Background(), &model.AssignCompanyRoleRequest{
		CustomerID: "cust-1",
		UserID:     "user-1",
		RoleName:   "admin",
	})
	assert.NoError(t, err)
}

func TestAssignCompanyRole_ValidationError(t *testing.T) {
	mockSvc := mocks.NewGrpcClient(t)
	uc := setupUseCase(t, mockSvc)

	err := uc.AssignCompanyRole(context.Background(), &model.AssignCompanyRoleRequest{
		CustomerID: "",
	})
	assert.Error(t, err)
}

// --- ListCompanyRoles ---

func TestListCompanyRoles_Success(t *testing.T) {
	mockSvc := mocks.NewGrpcClient(t)
	mockSvc.EXPECT().ListCompanyRoles(mock.Anything, mock.AnythingOfType("*customerv1.ListCompanyRolesRequest")).
		Return(&customerv1.ListCompanyRolesResponse{
			Roles: []*customerv1.CompanyRole{{Id: "r1", Name: "admin", CustomerId: "cust-1"}},
		}, nil)
	uc := setupUseCase(t, mockSvc)

	resp, err := uc.ListCompanyRoles(context.Background(), "cust-1")
	require.NoError(t, err)
	assert.Len(t, resp.Roles, 1)
	assert.Equal(t, "admin", resp.Roles[0].Name)
}

// --- GetUserPermissions ---

func TestGetUserPermissions_Success(t *testing.T) {
	mockSvc := mocks.NewGrpcClient(t)
	mockSvc.EXPECT().GetUserPermissions(mock.Anything, mock.AnythingOfType("*customerv1.GetUserPermissionsRequest")).
		Return(&customerv1.GetUserPermissionsResponse{
			Permissions: []string{"customer:read"},
		}, nil)
	uc := setupUseCase(t, mockSvc)

	resp, err := uc.GetUserPermissions(context.Background(), "user-1", "cust-1")
	require.NoError(t, err)
	assert.Equal(t, []string{"customer:read"}, resp.Permissions)
}
