package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/beabys/wms/inbound-service-bff/internal/api/v1"
	inboundv1 "github.com/beabys/wms/proto/gen/go/inbound/v1"
)

// mockGRPCClient implements GRPCInboundClient for testing.
type mockGRPCClient struct {
	createInboundResp   *inboundv1.Inbound
	createInboundErr    error
	getInboundResp      *inboundv1.Inbound
	getInboundErr       error
	listInboundsResp    []*inboundv1.Inbound
	listInboundsErr     error
	inspectInboundResp  *inboundv1.Inbound
	inspectInboundErr   error
	approveInboundResp  *inboundv1.Inbound
	approveInboundErr   error
	flagInboundResp     *inboundv1.Inbound
	flagInboundErr      error
	holdInboundResp     *inboundv1.Inbound
	holdInboundErr      error
	releaseInboundResp  *inboundv1.Inbound
	releaseInboundErr   error
}

func (m *mockGRPCClient) CreateInbound(_ context.Context, customerID, expectedDate, notes string, items []*inboundv1.InboundItem, token string) (*inboundv1.Inbound, error) {
	return m.createInboundResp, m.createInboundErr
}

func (m *mockGRPCClient) GetInbound(_ context.Context, id string, token string) (*inboundv1.Inbound, error) {
	return m.getInboundResp, m.getInboundErr
}

func (m *mockGRPCClient) ListInbounds(_ context.Context, customerID, status string, pageSize int32, pageToken string, token string) ([]*inboundv1.Inbound, error) {
	return m.listInboundsResp, m.listInboundsErr
}

func (m *mockGRPCClient) InspectInbound(_ context.Context, id, inspectorID, notes string, passed bool, token string) (*inboundv1.Inbound, error) {
	return m.inspectInboundResp, m.inspectInboundErr
}

func (m *mockGRPCClient) ApproveInbound(_ context.Context, id string, token string) (*inboundv1.Inbound, error) {
	return m.approveInboundResp, m.approveInboundErr
}

func (m *mockGRPCClient) FlagInbound(_ context.Context, id, reason string, token string) (*inboundv1.Inbound, error) {
	return m.flagInboundResp, m.flagInboundErr
}

func (m *mockGRPCClient) HoldInbound(_ context.Context, id, reason string, token string) (*inboundv1.Inbound, error) {
	return m.holdInboundResp, m.holdInboundErr
}

func (m *mockGRPCClient) ReleaseInbound(_ context.Context, id string, token string) (*inboundv1.Inbound, error) {
	return m.releaseInboundResp, m.releaseInboundErr
}

func (m *mockGRPCClient) Close() error {
	return nil
}

// setupTest creates an HttpServer with a mock gRPC client and chi router.
func setupTest(mock GRPCInboundClient) *chi.Mux {
	hs := NewHttpServer().
		SetLogger(zap.NewNop()).
		SetInboundClient(mock)

	r := chi.NewRouter()
	v1.HandlerWithOptions(hs, v1.ChiServerOptions{
		BaseRouter: r,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			w.WriteHeader(http.StatusInternalServerError)
		},
	})
	return r
}

// ---- SubmitInbound ----

func TestSubmitInbound_Success(t *testing.T) {
	mock := &mockGRPCClient{
		createInboundResp: &inboundv1.Inbound{
			Id:         "inbound-1",
			CustomerId: "cust-1",
			Status:     "pending",
		},
	}
	router := setupTest(mock)

	body := `{"customer_id":"cust-1","expected_date":"2024-01-15","items":[{"sku":"SKU-001","quantity_declared":10}]}`
	req := httptest.NewRequest("POST", "/v1/inbounds", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
	require.NotNil(t, apiResp["data"])
}

func TestSubmitInbound_InvalidBody(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/inbounds", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestSubmitInbound_GRPCError(t *testing.T) {
	mock := &mockGRPCClient{
		createInboundErr: grpc.ErrClientConnClosing,
	}
	router := setupTest(mock)

	body := `{"customer_id":"cust-1","expected_date":"2024-01-15","items":[{"sku":"SKU-001","quantity_declared":10}]}`
	req := httptest.NewRequest("POST", "/v1/inbounds", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
}

// ---- GetInboundQueue ----

func TestGetInboundQueue_Success(t *testing.T) {
	mock := &mockGRPCClient{
		listInboundsResp: []*inboundv1.Inbound{
			{Id: "inbound-1", CustomerId: "cust-1", Status: "pending"},
			{Id: "inbound-2", CustomerId: "cust-2", Status: "approved"},
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/inbounds/queue?page_size=10&status=pending", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
}

func TestGetInboundQueue_Empty(t *testing.T) {
	mock := &mockGRPCClient{
		listInboundsResp: []*inboundv1.Inbound{},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/inbounds/queue", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
}

func TestGetInboundQueue_GRPCError(t *testing.T) {
	mock := &mockGRPCClient{
		listInboundsErr: grpc.ErrClientConnClosing,
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/inbounds/queue", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
}

// ---- GetInbound ----

func TestGetInbound_Success(t *testing.T) {
	mock := &mockGRPCClient{
		getInboundResp: &inboundv1.Inbound{
			Id:         "inbound-1",
			CustomerId: "cust-1",
			Status:     "pending",
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/inbounds/inbound-1", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
}

func TestGetInbound_NotFound(t *testing.T) {
	mock := &mockGRPCClient{
		getInboundErr: status.Error(codes.NotFound, "inbound not found"),
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/inbounds/nonexistent", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNotFound, resp.Code)
}

func TestGetInbound_GRPCError(t *testing.T) {
	mock := &mockGRPCClient{
		getInboundErr: status.Error(codes.Internal, "internal server error"),
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/inbounds/inbound-1", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
}

// ---- ApproveInbound ----

func TestApproveInbound_Success(t *testing.T) {
	mock := &mockGRPCClient{
		approveInboundResp: &inboundv1.Inbound{
			Id:         "inbound-1",
			CustomerId: "cust-1",
			Status:     "approved",
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/inbounds/inbound-1/approve", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
}

func TestApproveInbound_GRPCError(t *testing.T) {
	mock := &mockGRPCClient{
		approveInboundErr: grpc.ErrClientConnClosing,
	}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/inbounds/inbound-1/approve", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
}

// ---- InspectInbound ----

func TestInspectInbound_Success(t *testing.T) {
	mock := &mockGRPCClient{
		inspectInboundResp: &inboundv1.Inbound{
			Id:         "inbound-1",
			CustomerId: "cust-1",
			Status:     "inspected",
		},
	}
	router := setupTest(mock)

	body := `{"inspector_id":"insp-1","passed":true,"notes":"all good"}`
	req := httptest.NewRequest("POST", "/v1/inbounds/inbound-1/inspect", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
}

func TestInspectInbound_InvalidBody(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/inbounds/inbound-1/inspect", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestInspectInbound_GRPCError(t *testing.T) {
	mock := &mockGRPCClient{
		inspectInboundErr: grpc.ErrClientConnClosing,
	}
	router := setupTest(mock)

	body := `{"inspector_id":"insp-1","passed":true}`
	req := httptest.NewRequest("POST", "/v1/inbounds/inbound-1/inspect", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
}

// ---- GetCustomerInbounds ----

func TestGetCustomerInbounds_Success(t *testing.T) {
	mock := &mockGRPCClient{
		listInboundsResp: []*inboundv1.Inbound{
			{Id: "inbound-1", CustomerId: "cust-1", Status: "pending"},
			{Id: "inbound-2", CustomerId: "cust-1", Status: "approved"},
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/inbounds/customer/cust-1?page_size=10", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
}

func TestGetCustomerInbounds_Empty(t *testing.T) {
	mock := &mockGRPCClient{
		listInboundsResp: []*inboundv1.Inbound{},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/inbounds/customer/cust-1", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var apiResp map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	require.True(t, apiResp["success"].(bool))
}

// ---- Missing Auth Token (Bug #3) ----

func TestSubmitInbound_MissingAuthToken(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	body := `{"customer_id":"cust-1","expected_date":"2024-01-15","items":[{"sku":"SKU-001","quantity_declared":10}]}`
	req := httptest.NewRequest("POST", "/v1/inbounds", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestGetInboundQueue_MissingAuthToken(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/inbounds/queue", nil)
	// No Authorization header
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestGetInbound_MissingAuthToken(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/inbounds/inbound-1", nil)
	// No Authorization header
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestInspectInbound_MissingAuthToken(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	body := `{"inspector_id":"insp-1","passed":true}`
	req := httptest.NewRequest("POST", "/v1/inbounds/inbound-1/inspect", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestApproveInbound_MissingAuthToken(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/inbounds/inbound-1/approve", nil)
	// No Authorization header
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestFlagInbound_MissingAuthToken(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	body := `{"reason":"damaged"}`
	req := httptest.NewRequest("POST", "/v1/inbounds/inbound-1/flag", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestHoldInbound_MissingAuthToken(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	body := `{"reason":"quality_check"}`
	req := httptest.NewRequest("POST", "/v1/inbounds/inbound-1/hold", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestReleaseInbound_MissingAuthToken(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/inbounds/inbound-1/release", nil)
	// No Authorization header
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestGetCustomerInbounds_MissingAuthToken(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/inbounds/customer/cust-1", nil)
	// No Authorization header
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

// ---- gRPC Status to HTTP Mapping (Bug #2) ----

func TestGrpcStatusToHTTP_FailedPrecondition(t *testing.T) {
	err := status.Error(codes.FailedPrecondition, "inbound not inspected yet")
	require.Equal(t, http.StatusBadRequest, grpcStatusToHTTP(err))
}

func TestGrpcStatusToHTTP_InvalidArgument(t *testing.T) {
	err := status.Error(codes.InvalidArgument, "invalid request")
	require.Equal(t, http.StatusBadRequest, grpcStatusToHTTP(err))
}

func TestGrpcStatusToHTTP_NotFound(t *testing.T) {
	err := status.Error(codes.NotFound, "inbound not found")
	require.Equal(t, http.StatusNotFound, grpcStatusToHTTP(err))
}

func TestGrpcStatusToHTTP_Unauthenticated(t *testing.T) {
	err := status.Error(codes.Unauthenticated, "unauthorized")
	require.Equal(t, http.StatusUnauthorized, grpcStatusToHTTP(err))
}

func TestGrpcStatusToHTTP_Internal(t *testing.T) {
	err := status.Error(codes.Internal, "internal error")
	require.Equal(t, http.StatusInternalServerError, grpcStatusToHTTP(err))
}

func TestGrpcStatusToHTTP_NonGRPCError(t *testing.T) {
	err := grpc.ErrClientConnClosing
	require.Equal(t, http.StatusInternalServerError, grpcStatusToHTTP(err))
}

func TestApproveInbound_FailedPrecondition(t *testing.T) {
	mock := &mockGRPCClient{
		approveInboundErr: status.Error(codes.FailedPrecondition, "inbound not inspected yet"),
	}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/inbounds/inbound-1/approve", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestApproveInbound_Unauthenticated(t *testing.T) {
	mock := &mockGRPCClient{
		approveInboundErr: status.Error(codes.Unauthenticated, "invalid token"),
	}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/inbounds/inbound-1/approve", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

// ---- NotFound ----

func TestNotFound(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/nonexistent", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNotFound, resp.Code)
}
