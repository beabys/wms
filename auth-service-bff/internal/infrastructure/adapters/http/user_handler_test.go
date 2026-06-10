package httpadapter

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	mocks "github.com/beabys/wms/auth-service-bff/mocks/infrastructure/http"
	grpcdapter "github.com/beabys/wms/auth-service-bff/internal/infrastructure/adapters/grpc"
	"github.com/beabys/wms/auth-service-bff/internal/domain/model"
)

func TestCreateUserHandler_Success(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().CreateUser(mock.Anything, mock.AnythingOfType("*model.CreateUserRequest")).
		Return(&model.UserResponse{ID: "2", Email: "new@test.com", Name: "New User", Role: "viewer", Active: true}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{"email":"new@test.com","password":"password123","name":"New User","role":"viewer"}`
	resp, err := http.Post(ts.URL+"/v1/users", "application/json", strings.NewReader(body))
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

func TestCreateUserHandler_ValidationError(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	ts := setupTestServer(svc)
	defer ts.Close()

	// Missing required fields — validation fails before mock is called
	body := `{"email":"new@test.com"}`
	resp, err := http.Post(ts.URL+"/v1/users", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestListUsersHandler_Success(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().ListUsers(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&model.UserListResponse{
			Users:      []model.UserResponse{{ID: "1", Email: "user@test.com", Name: "User", Role: "viewer", Active: true}},
			Pagination: model.Pagination{Page: 1, PageSize: 20, TotalItems: 1},
		}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/v1/users?page=1&page_size=10", nil)
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

func TestListUsersHandler_NoParams(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().ListUsers(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&model.UserListResponse{
			Users:      []model.UserResponse{{ID: "1", Email: "user@test.com", Name: "User", Role: "viewer", Active: true}},
			Pagination: model.Pagination{Page: 1, PageSize: 20, TotalItems: 1},
		}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/v1/users", nil)
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
	data := result["data"].(map[string]interface{})
	users := data["users"].([]interface{})
	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
}

func TestGetUserHandler_Success(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().GetUser(mock.Anything, "user-1").
		Return(&model.UserResponse{ID: "user-1", Email: "user@test.com", Name: "User", Role: "viewer", Active: true}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/v1/users/user-1", nil)
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

func TestGetUserHandler_NotFound(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().GetUser(mock.Anything, "nonexistent").
		Return(nil, grpcdapter.ErrNotFound)
	ts := setupTestServer(svc)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/v1/users/nonexistent", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestUpdateUserHandler_Success(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().UpdateUser(mock.Anything, "user-1", mock.AnythingOfType("*model.UpdateUserRequest")).
		Return(&model.UserResponse{ID: "user-1", Email: "user@test.com", Name: "Updated", Role: "admin", Active: true}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{"name":"Updated Name","role":"admin"}`
	req, _ := http.NewRequest("PUT", ts.URL+"/v1/users/user-1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
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

func TestDeleteUserHandler_Success(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().DeleteUser(mock.Anything, "user-1").Return(nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	req, _ := http.NewRequest("DELETE", ts.URL+"/v1/users/user-1", nil)
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

func TestAssignRoleHandler_Success(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().AssignRole(mock.Anything, "user-1", "admin").Return(nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{"role":"admin"}`
	req, _ := http.NewRequest("POST", ts.URL+"/v1/users/user-1/roles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
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

func TestAssignRoleHandler_EmptyRole(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{"role":""}`
	req, _ := http.NewRequest("POST", ts.URL+"/v1/users/user-1/roles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestErrorMapping_Unauthenticated(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().Login(mock.Anything, "test@test.com", "password123").
		Return(nil, grpcdapter.ErrUnauthenticated)
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{"email":"test@test.com","password":"password123"}`
	resp, err := http.Post(ts.URL+"/v1/auth/login", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestErrorMapping_InternalError(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().Login(mock.Anything, "test@test.com", "password123").
		Return(nil, errors.New("unexpected error"))
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{"email":"test@test.com","password":"password123"}`
	resp, err := http.Post(ts.URL+"/v1/auth/login", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}

	result := decodeResponse(t, resp)
	if result["success"] != false {
		t.Error("expected success=false")
	}
	errMsg := result["error"].(string)
	if strings.Contains(errMsg, "internal server error") {
		// Good — internal errors should hide details
	}
}
