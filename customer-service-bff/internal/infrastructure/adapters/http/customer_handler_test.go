package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	v1 "github.com/beabys/wms/customer-service-bff/internal/api/v1"
	"github.com/beabys/wms/customer-service-bff/internal/application"
)

// mockCustomerClient implements application.CustomerService for testing.
type mockCustomerClient struct {
	createCustomerResp  *application.Customer
	createCustomerErr   error
	getCustomerResp     *application.Customer
	getCustomerErr      error
	listCustomersResp   *application.ListResult
	listCustomersErr    error
	approveCustomerResp *application.Customer
	approveCustomerErr  error
	suspendCustomerResp *application.Customer
	suspendCustomerErr  error
	inviteCustomerResp  string
	inviteCustomerErr   error
}

func (m *mockCustomerClient) CreateCustomer(_ context.Context, _ application.CreateCustomerRequest, _ string) (*application.Customer, error) {
	return m.createCustomerResp, m.createCustomerErr
}

func (m *mockCustomerClient) GetCustomer(_ context.Context, _, _ string) (*application.Customer, error) {
	return m.getCustomerResp, m.getCustomerErr
}

func (m *mockCustomerClient) ListCustomers(_ context.Context, _ string, _, _ int32, _ string) (*application.ListResult, error) {
	return m.listCustomersResp, m.listCustomersErr
}

func (m *mockCustomerClient) ApproveCustomer(_ context.Context, _, _ string) (*application.Customer, error) {
	return m.approveCustomerResp, m.approveCustomerErr
}

func (m *mockCustomerClient) SuspendCustomer(_ context.Context, _, _, _ string) (*application.Customer, error) {
	return m.suspendCustomerResp, m.suspendCustomerErr
}

func (m *mockCustomerClient) InviteCustomer(_ context.Context, _, _ string) (string, error) {
	return m.inviteCustomerResp, m.inviteCustomerErr
}

// setupTest creates an HttpServer with a mock customer client and chi router.
func setupTest(mock application.CustomerService) *chi.Mux {
	hs := NewHttpServer(&Config{}, zap.NewNop(), mock)
	r := chi.NewRouter()
	v1.HandlerWithOptions(hs, v1.ChiServerOptions{
		BaseRouter: r,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			w.WriteHeader(http.StatusInternalServerError)
		},
	})
	return r
}

// decodeResponse unmarshals JSON response body into a map.
func decodeResponse(t *testing.T, resp *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	return result
}

// ---- InviteCustomer ----

func TestInviteCustomer_Success(t *testing.T) {
	mock := &mockCustomerClient{
		inviteCustomerResp: "invite-token-abc",
	}
	router := setupTest(mock)

	body := `{"email":"test@example.com"}`
	req := httptest.NewRequest("POST", "/v1/customers/invite", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	result := decodeResponse(t, resp)
	require.True(t, result["success"].(bool))
	require.NotNil(t, result["data"])
}

func TestInviteCustomer_InvalidBody(t *testing.T) {
	mock := &mockCustomerClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/customers/invite", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestInviteCustomer_GRPCError(t *testing.T) {
	mock := &mockCustomerClient{
		inviteCustomerErr: grpc.ErrClientConnClosing,
	}
	router := setupTest(mock)

	body := `{"email":"test@example.com"}`
	req := httptest.NewRequest("POST", "/v1/customers/invite", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
}

// ---- RegisterCustomer ----

func TestRegisterCustomer_Success(t *testing.T) {
	mock := &mockCustomerClient{
		createCustomerResp: &application.Customer{
			ID:          "cust-1",
			CompanyName: "Test Corp",
			Status:      "pending",
		},
	}
	router := setupTest(mock)

	body := `{"company_name":"Test Corp","vat_number":"VAT123","line1":"123 Main St","city":"Springfield","country":"US"}`
	req := httptest.NewRequest("POST", "/v1/customers/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	result := decodeResponse(t, resp)
	require.True(t, result["success"].(bool))
	require.NotNil(t, result["data"])
}

func TestRegisterCustomer_InvalidBody(t *testing.T) {
	mock := &mockCustomerClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/customers/register", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestRegisterCustomer_GRPCError(t *testing.T) {
	mock := &mockCustomerClient{
		createCustomerErr: grpc.ErrClientConnClosing,
	}
	router := setupTest(mock)

	body := `{"company_name":"Test Corp","vat_number":"VAT123","line1":"123 Main St","city":"Springfield","country":"US"}`
	req := httptest.NewRequest("POST", "/v1/customers/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
}

// ---- GetCustomerApprovalQueue ----

func TestGetCustomerApprovalQueue_Success(t *testing.T) {
	mock := &mockCustomerClient{
		listCustomersResp: &application.ListResult{
			Customers: []application.Customer{
				{ID: "cust-1", CompanyName: "Test Corp", Status: "pending"},
			},
			Total:    1,
			Page:     1,
			PageSize: 20,
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/customers/queue?page=1&page_size=20", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	result := decodeResponse(t, resp)
	require.True(t, result["success"].(bool))
	require.NotNil(t, result["data"])
}

func TestGetCustomerApprovalQueue_GRPCError(t *testing.T) {
	mock := &mockCustomerClient{
		listCustomersErr: errors.New("service unavailable"),
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/customers/queue", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
}

// ---- GetCustomer ----

func TestGetCustomer_Success(t *testing.T) {
	mock := &mockCustomerClient{
		getCustomerResp: &application.Customer{
			ID:          "cust-1",
			CompanyName: "Test Corp",
			Status:      "active",
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/customers/cust-1", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	result := decodeResponse(t, resp)
	require.True(t, result["success"].(bool))
	require.NotNil(t, result["data"])
}

func TestGetCustomer_NotFound(t *testing.T) {
	mock := &mockCustomerClient{
		getCustomerErr: errors.New("customer not found"),
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/customers/nonexistent", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNotFound, resp.Code)
}

// ---- ApproveCustomer ----

func TestApproveCustomer_Success(t *testing.T) {
	mock := &mockCustomerClient{
		approveCustomerResp: &application.Customer{
			ID:          "cust-1",
			CompanyName: "Test Corp",
			Status:      "active",
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/customers/cust-1/approve", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	result := decodeResponse(t, resp)
	require.True(t, result["success"].(bool))
	require.NotNil(t, result["data"])
}

func TestApproveCustomer_GRPCError(t *testing.T) {
	mock := &mockCustomerClient{
		approveCustomerErr: errors.New("approval failed"),
	}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/customers/cust-1/approve", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
}

// ---- SuspendCustomer ----

func TestSuspendCustomer_Success(t *testing.T) {
	mock := &mockCustomerClient{
		suspendCustomerResp: &application.Customer{
			ID:          "cust-1",
			CompanyName: "Test Corp",
			Status:      "suspended",
		},
	}
	router := setupTest(mock)

	body := `{"reason":"Violated terms of service"}`
	req := httptest.NewRequest("POST", "/v1/customers/cust-1/suspend", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	result := decodeResponse(t, resp)
	require.True(t, result["success"].(bool))
	require.NotNil(t, result["data"])
}

func TestSuspendCustomer_NoBody(t *testing.T) {
	// SuspendCustomer treats decode failure as non-fatal (reason is optional).
	mock := &mockCustomerClient{
		suspendCustomerResp: &application.Customer{
			ID:          "cust-1",
			CompanyName: "Test Corp",
			Status:      "suspended",
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/customers/cust-1/suspend", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	result := decodeResponse(t, resp)
	require.True(t, result["success"].(bool))
}

func TestSuspendCustomer_GRPCError(t *testing.T) {
	mock := &mockCustomerClient{
		suspendCustomerErr: errors.New("suspend failed"),
	}
	router := setupTest(mock)

	body := `{"reason":"Violated terms"}`
	req := httptest.NewRequest("POST", "/v1/customers/cust-1/suspend", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
}

// ---- GetCustomerDashboard ----

func TestGetCustomerDashboard_Success(t *testing.T) {
	mock := &mockCustomerClient{
		getCustomerResp: &application.Customer{
			ID:          "cust-1",
			CompanyName: "Test Corp",
			Status:      "active",
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/customers/cust-1/dashboard", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	result := decodeResponse(t, resp)
	require.True(t, result["success"].(bool))
	require.NotNil(t, result["data"])
}

func TestGetCustomerDashboard_NotFound(t *testing.T) {
	mock := &mockCustomerClient{
		getCustomerErr: errors.New("customer not found"),
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/customers/cust-999/dashboard", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNotFound, resp.Code)
}

// ---- NotFound ----

func TestNotFound(t *testing.T) {
	mock := &mockCustomerClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/nonexistent", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNotFound, resp.Code)
}
