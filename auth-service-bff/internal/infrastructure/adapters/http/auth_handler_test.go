package httpadapter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	mocks "github.com/beabys/wms/auth-service-bff/mocks/infrastructure/http"
	grpcdapter "github.com/beabys/wms/auth-service-bff/internal/infrastructure/adapters/grpc"
	"github.com/beabys/wms/auth-service-bff/internal/domain/model"
	"github.com/beabys/wms/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func setupTestServer(svc *mocks.AuthHandlerService) *httptest.Server {
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

// --- Auth handler tests ---

func TestLoginHandler_Success(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().Login(mock.Anything, "test@example.com", "password123").
		Return(&model.AuthResponse{AccessToken: "mock-token", RefreshToken: "mock-refresh", ExpiresIn: 3600}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{"email":"test@example.com","password":"password123"}`
	resp, err := http.Post(ts.URL+"/v1/auth/login", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	result := decodeResponse(t, resp)
	if result["success"] != true {
		t.Error("expected success=true")
	}
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data object")
	}
	if data["access_token"] != "mock-token" {
		t.Errorf("expected mock-token, got %v", data["access_token"])
	}
}

func TestLoginHandler_ValidationError(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	ts := setupTestServer(svc)
	defer ts.Close()

	// Missing password — validation fails before mock is called
	body := `{"email":"test@example.com"}`
	resp, err := http.Post(ts.URL+"/v1/auth/login", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}

	result := decodeResponse(t, resp)
	if result["success"] != false {
		t.Error("expected success=false")
	}
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	ts := setupTestServer(svc)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/v1/auth/login", "application/json", strings.NewReader(`{invalid`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestRefreshTokenHandler_Success(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().RefreshToken(mock.Anything, "valid-refresh-token").
		Return(&model.AuthResponse{AccessToken: "new-token", RefreshToken: "new-refresh", ExpiresIn: 3600}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{"refresh_token":"valid-refresh-token"}`
	resp, err := http.Post(ts.URL+"/v1/auth/refresh", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	result := decodeResponse(t, resp)
	if result["success"] != true {
		t.Error("expected success=true")
	}
}

func TestRefreshTokenHandler_EmptyToken(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{"refresh_token":""}`
	resp, err := http.Post(ts.URL+"/v1/auth/refresh", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestGetMeHandler_Success(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().GetMe(mock.Anything).
		Return(&model.UserResponse{ID: "1", Email: "test@example.com", Name: "Test", Role: "admin", Active: true}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	result := decodeResponse(t, resp)
	if result["success"] != true {
		t.Error("expected success=true")
	}
}

// --- Invite handler tests ---

func TestInviteUserHandler_Success(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().InviteUser(mock.Anything, "invite@example.com", "").
		Return(&model.InviteUserResponse{
			Token:      "mock-invite-token",
			InviteLink: "https://portal.example.com/register?token=mock-invite-token",
			ExpiresAt:  1800000000,
		}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{"email":"invite@example.com"}`
	req, _ := http.NewRequest("POST", ts.URL+"/v1/auth/invite", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	result := decodeResponse(t, resp)
	if result["success"] != true {
		t.Error("expected success=true")
	}
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data object")
	}
	if data["token"] != "mock-invite-token" {
		t.Errorf("expected mock-invite-token, got %v", data["token"])
	}
}

func TestInviteUserHandler_MissingEmail(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{}`
	req, _ := http.NewRequest("POST", ts.URL+"/v1/auth/invite", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestInviteUserHandler_InvalidJSON(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	ts := setupTestServer(svc)
	defer ts.Close()

	req, _ := http.NewRequest("POST", ts.URL+"/v1/auth/invite", strings.NewReader(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

// --- Logout handler tests ---

func TestLogoutHandler_Success(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().Logout(mock.Anything, "some-refresh-token").Return(nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{"refresh_token":"some-refresh-token"}`
	req, _ := http.NewRequest("POST", ts.URL+"/v1/auth/logout", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	result := decodeResponse(t, resp)
	if result["success"] != true {
		t.Error("expected success=true")
	}
}

func TestLogoutHandler_NoBody(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().Logout(mock.Anything, "").Return(nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	req, _ := http.NewRequest("POST", ts.URL+"/v1/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	result := decodeResponse(t, resp)
	if result["success"] != true {
		t.Error("expected success=true")
	}
}

func TestLogoutHandler_EmptyBody(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().Logout(mock.Anything, "").Return(nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{}`
	req, _ := http.NewRequest("POST", ts.URL+"/v1/auth/logout", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestLogoutHandler_ServiceError(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().Logout(mock.Anything, "some-token").Return(grpcdapter.ErrInternal)
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{"refresh_token":"some-token"}`
	req, _ := http.NewRequest("POST", ts.URL+"/v1/auth/logout", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

// --- ListInvites handler tests ---

func TestListInvitesHandler_Success(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().ListInvites(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&model.InviteListResponse{
			Invites: []model.InviteEntry{
				{ID: "1", Email: "invited@example.com", Token: "abc123", InvitedBy: "admin", Status: "pending", ExpiresAt: 1900000000, CreatedAt: 1800000000},
			},
			Pagination: model.Pagination{Page: 1, PageSize: 20, TotalItems: 1},
		}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/v1/auth/invites?page=1&page_size=10", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	result := decodeResponse(t, resp)
	if result["success"] != true {
		t.Error("expected success=true")
	}
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data object")
	}
	invites := data["invites"].([]interface{})
	if len(invites) != 1 {
		t.Errorf("expected 1 invite, got %d", len(invites))
	}
}

func TestListInvitesHandler_WithFilters(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().ListInvites(
		mock.Anything,
		mock.MatchedBy(func(status *string) bool { return status != nil && *status == "pending" }),
		mock.MatchedBy(func(expired *bool) bool { return expired != nil && !*expired }),
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(page *int) bool { return page != nil && *page == 1 }),
		mock.Anything,
	).
		Return(&model.InviteListResponse{
			Invites:    []model.InviteEntry{{ID: "1", Email: "filtered@example.com", Token: "xyz", Status: "pending"}},
			Pagination: model.Pagination{Page: 1, PageSize: 20, TotalItems: 1},
		}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/v1/auth/invites?status=pending&expired=false&page=1", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestListInvitesHandler_NoAuth(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().ListInvites(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&model.InviteListResponse{
			Invites:    []model.InviteEntry{{ID: "1", Email: "invited@example.com", Token: "abc123", InvitedBy: "admin", Status: "pending"}},
			Pagination: model.Pagination{Page: 1, PageSize: 20, TotalItems: 1},
		}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	// No Authorization header
	resp, err := http.Get(ts.URL + "/v1/auth/invites")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Should still work — handler doesn't enforce auth, middleware does
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestListInvitesHandler_ServiceError(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().ListInvites(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, grpcdapter.ErrInternal)
	ts := setupTestServer(svc)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/v1/auth/invites", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

// --- CancelInvite handler tests ---

func TestCancelInviteHandler_Success(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().CancelInvite(mock.Anything, "test-token-123").
		Return(&model.CancelInviteResponse{Success: true}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	req, _ := http.NewRequest("POST", ts.URL+"/v1/auth/invites/test-token-123/cancel", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	result := decodeResponse(t, resp)
	if result["success"] != true {
		t.Error("expected success=true")
	}
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data object")
	}
	if data["success"] != true {
		t.Errorf("expected success=true, got %v", data["success"])
	}
}

func TestCancelInviteHandler_ServiceError(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().CancelInvite(mock.Anything, "some-token").
		Return(nil, grpcdapter.ErrInternal)
	ts := setupTestServer(svc)
	defer ts.Close()

	req, _ := http.NewRequest("POST", ts.URL+"/v1/auth/invites/some-token/cancel", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

func TestCancelInviteHandler_NoAuth(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().CancelInvite(mock.Anything, "test-token").Return(&model.CancelInviteResponse{Success: true}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/v1/auth/invites/test-token/cancel", "application/json", nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Should still work — handler doesn't enforce auth, middleware does
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestHealthEndpoint(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	ts := setupTestServer(svc)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
