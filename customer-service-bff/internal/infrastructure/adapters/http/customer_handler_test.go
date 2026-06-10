package httpadapter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	mocks "github.com/beabys/wms/customer-service-bff/mocks/infrastructure/http"
	grpcdapter "github.com/beabys/wms/customer-service-bff/internal/infrastructure/adapters/grpc"
	"github.com/beabys/wms/customer-service-bff/internal/domain/model"
	"github.com/beabys/wms/pkg/logger"
)

func setupTestServer(t *testing.T, svc *mocks.CustomerHandlerService) *httptest.Server {
	t.Helper()
	log, _ := logger.NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	srv := NewServer(svc, log)
	handler, _ := NewMuxHandler(srv, log, []string{"*"})
	return httptest.NewServer(handler)
}

func decodeResponse(t *testing.T, resp *http.Response) map[string]interface{} {
	t.Helper()
	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	return body
}

// --- RegisterCustomer tests ---

func TestRegisterCustomerHandler_Success(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().RegisterCustomer(mock.Anything, mock.AnythingOfType("*model.RegisterCustomerRequest")).
		Return(&model.RegisterCustomerResponse{
			Customer:    &model.CustomerResponse{ID: "mock-cust", CompanyName: "ACME Corp", Status: "pending"},
			AccessToken: "mock-token",
		}, nil)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	body := `{"token":"invite-token","company_name":"ACME Corp","email":"admin@acme.com","password":"secure-pass"}`
	resp, err := http.Post(ts.URL+"/v1/customers/register", "application/json", strings.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.True(t, result["success"].(bool))
}

func TestRegisterCustomerHandler_ValidationError(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().RegisterCustomer(mock.Anything, mock.AnythingOfType("*model.RegisterCustomerRequest")).
		Return(nil, grpcdapter.ErrInvalidArgument)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	body := `{"token":"invite-token"}`
	resp, err := http.Post(ts.URL+"/v1/customers/register", "application/json", strings.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.False(t, result["success"].(bool))
}

func TestRegisterCustomerHandler_InvalidJSON(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/v1/customers/register", "application/json", strings.NewReader(`{invalid`))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// --- ListCustomers tests ---

func TestListCustomersHandler_Success(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().ListCustomers(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&model.CustomerListResponse{
			Customers:  []model.CustomerResponse{{ID: "c1", CompanyName: "Co 1"}},
			Pagination: model.Pagination{Page: 1, PageSize: 20, TotalItems: 1},
		}, nil)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/customers?page=1&page_size=10", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.True(t, result["success"].(bool))
}

func TestListCustomersHandler_NoAuth(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().ListCustomers(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&model.CustomerListResponse{
			Customers:  []model.CustomerResponse{{ID: "c1", CompanyName: "Co 1"}},
			Pagination: model.Pagination{Page: 1, PageSize: 20, TotalItems: 1},
		}, nil)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	// No auth header — should still work (JWT extractor just puts empty token in context)
	resp, err := http.Get(ts.URL + "/v1/customers")
	require.NoError(t, err)
	defer resp.Body.Close()

	// This should work since the handler doesn't check JWT itself
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// --- GetMyCustomer tests ---

func TestGetMyCustomerHandler_Success(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().GetMyCustomer(mock.Anything, "test-token").
		Return(&model.CustomerResponse{ID: "mock-cust", CompanyName: "Mock Co", Status: "active", CompanyAdminID: "test-token"}, nil)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/customers/me", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.True(t, result["success"].(bool))
	data, ok := result["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "mock-cust", data["id"])
	assert.Equal(t, "active", data["status"])
}

func TestGetMyCustomerHandler_NotFound(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().GetMyCustomer(mock.Anything, "test-token").
		Return(nil, grpcdapter.ErrNotFound)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/customers/me", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.False(t, result["success"].(bool))
}

func TestGetMyCustomerHandler_NoAuth(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	// No auth header — should return 401 (handled by the handler's JWT check)
	resp, err := http.Get(ts.URL + "/v1/customers/me")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.False(t, result["success"].(bool))
}

func TestGetMyCustomerHandler_GrpcNotFound(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().GetMyCustomer(mock.Anything, "test-token").
		Return(nil, grpcdapter.ErrNotFound)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/customers/me", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetMyCustomerHandler_InternalError(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().GetMyCustomer(mock.Anything, "test-token").
		Return(nil, grpcdapter.ErrInternal)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/customers/me", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// --- GetCustomer tests ---

func TestGetCustomerHandler_Success(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().GetCustomer(mock.Anything, "cust-123").
		Return(&model.CustomerResponse{ID: "cust-123", CompanyName: "Mock Co", Status: "active"}, nil)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/customers/cust-123", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.True(t, result["success"].(bool))
}

// --- UpdateCustomer tests ---

func TestUpdateCustomerHandler_Success(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().UpdateCustomer(mock.Anything, "cust-123", mock.AnythingOfType("*model.UpdateCustomerRequest")).
		Return(&model.CustomerResponse{ID: "cust-123", City: "New York", Status: "active"}, nil)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	body := `{"phone":"555-0100","city":"New York"}`
	req, _ := http.NewRequest(http.MethodPut, ts.URL+"/v1/customers/cust-123", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.True(t, result["success"].(bool))
}

func TestUpdateCustomerHandler_InvalidJSON(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPut, ts.URL+"/v1/customers/cust-123", strings.NewReader(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// --- ApproveCustomer tests ---

func TestApproveCustomerHandler_Success(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().ApproveCustomer(mock.Anything, "cust-123", "admin-token").
		Return(&model.CustomerResponse{ID: "cust-123", Status: "active"}, nil)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/customers/cust-123/approve", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.True(t, result["success"].(bool))
}

// --- RejectCustomer tests ---

func TestRejectCustomerHandler_Success(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().RejectCustomer(mock.Anything, "cust-123", "invalid documents").Return(nil)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	body := `{"reason":"invalid documents"}`
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/customers/cust-123/reject", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.True(t, result["success"].(bool))
}

// --- SuspendCustomer tests ---

func TestSuspendCustomerHandler_Success(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().SuspendCustomer(mock.Anything, "cust-123", "policy violation").Return(nil)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	body := `{"reason":"policy violation"}`
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/customers/cust-123/suspend", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// --- RestoreCustomer tests ---

func TestRestoreCustomerHandler_Success(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().RestoreCustomer(mock.Anything, "cust-123", "admin-token").
		Return(&model.CustomerResponse{ID: "cust-123", Status: "active"}, nil)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/customers/cust-123/restore", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.True(t, result["success"].(bool))
}

func TestRestoreCustomerHandler_NotFound(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().RestoreCustomer(mock.Anything, "cust-999", "admin-token").
		Return(nil, grpcdapter.ErrNotFound)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/customers/cust-999/restore", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestRestoreCustomerHandler_NotSuspended(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().RestoreCustomer(mock.Anything, "cust-123", "admin-token").
		Return(nil, grpcdapter.ErrInvalidArgument)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/customers/cust-123/restore", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRestoreCustomerHandler_InternalError(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().RestoreCustomer(mock.Anything, "cust-123", "admin-token").
		Return(nil, grpcdapter.ErrInternal)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/customers/cust-123/restore", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// --- ListAuditLogs tests ---

func TestListAuditLogsHandler_Success(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().ListAuditLogs(mock.Anything, "cust-123", 1, 10).
		Return(&model.AuditLogListResponse{
			Entries: []model.AuditEntry{
				{ID: "audit-1", CustomerID: "cust-123", Action: "approved", PerformedBy: "admin-1", CreatedAt: "2024-01-15T10:00:00Z"},
			},
			Pagination: model.Pagination{Page: 1, PageSize: 20, TotalItems: 1},
		}, nil)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/customers/cust-123/audit?page=1&page_size=10", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.True(t, result["success"].(bool))
}

func TestListAuditLogsHandler_NotFound(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().ListAuditLogs(mock.Anything, "nonexistent", 1, 20).
		Return(nil, grpcdapter.ErrNotFound)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/customers/nonexistent/audit", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestListAuditLogsHandler_InternalError(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().ListAuditLogs(mock.Anything, "cust-123", 1, 20).
		Return(nil, grpcdapter.ErrInternal)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/customers/cust-123/audit", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// --- AssignCompanyRole tests ---

func TestAssignCompanyRoleHandler_Success(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().AssignCompanyRole(mock.Anything, mock.AnythingOfType("*model.AssignCompanyRoleRequest")).Return(nil)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	body := `{"customer_id":"cust-1","user_id":"user-1","role_name":"admin"}`
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/companies/roles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAssignCompanyRoleHandler_ValidationError(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().AssignCompanyRole(mock.Anything, mock.AnythingOfType("*model.AssignCompanyRoleRequest")).
		Return(grpcdapter.ErrInvalidArgument)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	body := `{"customer_id":"cust-1"}` // missing user_id and role_name
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/companies/roles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// --- ListCompanyRoles tests ---

func TestListCompanyRolesHandler_Success(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().ListCompanyRoles(mock.Anything, "cust-1").
		Return(&model.CompanyRoleListResponse{
			Roles: []model.CompanyRoleResponse{{ID: "r1", Name: "admin"}},
		}, nil)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/companies/roles?customer_id=cust-1", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.True(t, result["success"].(bool))
}

func TestListCompanyRolesHandler_MissingCustomerID(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/companies/roles", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// --- GetUserPermissions tests ---

func TestGetUserPermissionsHandler_Success(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().GetUserPermissions(mock.Anything, "user-1", "cust-1").
		Return(&model.PermissionsResponse{Permissions: []string{"read"}}, nil)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/companies/permissions?user_id=user-1&customer_id=cust-1", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.True(t, result["success"].(bool))
}

func TestGetUserPermissionsHandler_MissingParams(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/companies/permissions", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// --- Health endpoint tests ---

func TestHealthEndpoint(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	ts := setupTestServer(t, svc)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.True(t, result["success"].(bool))
	assert.Equal(t, "ok", result["data"].(map[string]interface{})["status"])
}

// --- Error handler tests ---

func TestHandlerGrpcErrorMapping(t *testing.T) {
	svc := mocks.NewCustomerHandlerService(t)
	svc.EXPECT().GetCustomer(mock.Anything, "cust-999").
		Return(nil, fmt.Errorf("some error"))
	ts := setupTestServer(t, svc)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/customers/cust-999", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	result := decodeResponse(t, resp)
	assert.False(t, result["success"].(bool))
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	// Non-sentinel errors keep their message
	assert.Equal(t, "some error", result["error"])
}
