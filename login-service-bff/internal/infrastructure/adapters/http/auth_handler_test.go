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

	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
)

// mockGRPCClient implements GRPCAuthClient for testing.
type mockGRPCClient struct {
	loginResp        *authv1.LoginResponse
	loginErr         error
	registerResp     *authv1.CreateUserResponse
	registerErr      error
	refreshResp      *authv1.RefreshTokenResponse
	refreshErr       error
	validateResp     *authv1.ValidateTokenResponse
	validateErr      error
	getUserResp      *authv1.GetUserResponse
	getUserErr       error
	deleteErr        error
}

func (m *mockGRPCClient) Login(_ context.Context, email, password string) (*authv1.LoginResponse, error) {
	return m.loginResp, m.loginErr
}

func (m *mockGRPCClient) Register(_ context.Context, email, password, companyName string) (*authv1.CreateUserResponse, error) {
	return m.registerResp, m.registerErr
}

func (m *mockGRPCClient) RefreshToken(_ context.Context, refreshToken string) (*authv1.RefreshTokenResponse, error) {
	return m.refreshResp, m.refreshErr
}

func (m *mockGRPCClient) ValidateToken(_ context.Context, token string) (*authv1.ValidateTokenResponse, error) {
	return m.validateResp, m.validateErr
}

func (m *mockGRPCClient) GetUser(_ context.Context, userID string, _ string) (*authv1.GetUserResponse, error) {
	return m.getUserResp, m.getUserErr
}

func (m *mockGRPCClient) DeleteRefreshToken(_ context.Context, token string) error {
	return m.deleteErr
}

func (m *mockGRPCClient) Close() error {
	return nil
}

// setupTest creates an HttpServer with a mock gRPC client and chi router.
func setupTest(mock GRPCAuthClient) *chi.Mux {
	hs := NewHttpServer().
		SetLogger(zap.NewNop()).
		SetGRPCClient(mock)

	r := chi.NewRouter()
	r.Post("/v1/auth/login", hs.Postv1AuthLogin)
	r.Post("/v1/auth/register", hs.Postv1AuthRegister)
	r.Post("/v1/auth/refresh", hs.Postv1AuthRefresh)
	r.Post("/v1/auth/logout", hs.Postv1AuthLogout)
	r.Get("/v1/auth/me", hs.Getv1AuthMe)
	return r
}

func TestLogin_Success(t *testing.T) {
	mock := &mockGRPCClient{
		loginResp: &authv1.LoginResponse{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			User: &authv1.User{
				Id:    "user-1",
				Email: "test@example.com",
				Role:  "customer",
			},
		},
	}

	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/auth/login",
		strings.NewReader(`{"email":"test@example.com","password":"pass123"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var body map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.True(t, body["success"].(bool))
	data := body["data"].(map[string]interface{})
	require.Equal(t, "test-access-token", data["access_token"])
	require.Equal(t, "test-refresh-token", data["refresh_token"])
}

func TestLogin_InvalidBody(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/auth/login",
		strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestLogin_Unauthorized(t *testing.T) {
	mock := &mockGRPCClient{
		loginErr: grpc.ErrClientConnClosing,
	}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/auth/login",
		strings.NewReader(`{"email":"test@example.com","password":"wrong"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestRegister_Success(t *testing.T) {
	mock := &mockGRPCClient{
		registerResp: &authv1.CreateUserResponse{
			User: &authv1.User{
				Id:    "new-user-1",
				Email: "new@example.com",
				Role:  "customer",
			},
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/auth/register",
		strings.NewReader(`{"email":"new@example.com","password":"pass123"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var body map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.True(t, body["success"].(bool))
}

func TestRefresh_Success(t *testing.T) {
	mock := &mockGRPCClient{
		refreshResp: &authv1.RefreshTokenResponse{
			AccessToken:  "new-access-token",
			RefreshToken: "new-refresh-token",
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/auth/refresh",
		strings.NewReader(`{"refresh_token":"old-token"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var body map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.True(t, body["success"].(bool))
}

func TestMe_Success(t *testing.T) {
	mock := &mockGRPCClient{
		validateResp: &authv1.ValidateTokenResponse{
			Valid:  true,
			UserId: "user-1",
			Role:   "customer",
		},
		getUserResp: &authv1.GetUserResponse{
			User: &authv1.User{
				Id:    "user-1",
				Email: "user@example.com",
				Role:  "customer",
			},
		},
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var body map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.True(t, body["success"].(bool))
}

func TestMe_Unauthorized(t *testing.T) {
	mock := &mockGRPCClient{
		validateErr: grpc.ErrClientConnClosing,
	}
	router := setupTest(mock)

	req := httptest.NewRequest("GET", "/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestLogout_Success(t *testing.T) {
	mock := &mockGRPCClient{}
	router := setupTest(mock)

	req := httptest.NewRequest("POST", "/v1/auth/logout",
		strings.NewReader(`{"refresh_token":"some-token"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var body map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.True(t, body["success"].(bool))
}
