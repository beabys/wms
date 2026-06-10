package httpadapter

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	mocks "github.com/beabys/wms/auth-service-bff/mocks/infrastructure/http"
	"github.com/beabys/wms/auth-service-bff/internal/domain/model"
)

func TestListRolesHandler_Success(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().ListRoles(mock.Anything).
		Return(&model.RoleListResponse{Roles: []model.RoleResponse{{ID: "1", Name: "admin", Description: "Admin role"}}}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/v1/roles", nil)
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

	data := result["data"].(map[string]interface{})
	roles := data["roles"].([]interface{})
	if len(roles) != 1 {
		t.Errorf("expected 1 role, got %d", len(roles))
	}

	firstRole := roles[0].(map[string]interface{})
	if firstRole["name"] != "admin" {
		t.Errorf("expected role name 'admin', got %v", firstRole["name"])
	}
}

func TestCreateRoleHandler_Success(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().CreateRole(mock.Anything, mock.AnythingOfType("*model.CreateRoleRequest")).
		Return(&model.RoleResponse{ID: "1", Name: "editor", Description: "Can edit content", IsSystem: false}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{"name":"editor","description":"Can edit content"}`
	req, _ := http.NewRequest("POST", ts.URL+"/v1/roles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Expected status: 201 Created
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	result := decodeResponse(t, resp)
	if result["success"] != true {
		t.Error("expected success=true")
	}

	data := result["data"].(map[string]interface{})
	if data["name"] != "editor" {
		t.Errorf("expected name 'editor', got %v", data["name"])
	}
}

func TestCreateRoleHandler_MissingName(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{"description":"Missing name"}`
	req, _ := http.NewRequest("POST", ts.URL+"/v1/roles", strings.NewReader(body))
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

func TestCreateRoleHandler_WithPermissions(t *testing.T) {
	svc := mocks.NewAuthHandlerService(t)
	svc.EXPECT().CreateRole(mock.Anything, mock.AnythingOfType("*model.CreateRoleRequest")).
		Return(&model.RoleResponse{ID: "1", Name: "custom", Description: "Custom role", IsSystem: false}, nil)
	ts := setupTestServer(svc)
	defer ts.Close()

	body := `{
		"name":"custom",
		"description":"Custom role",
		"permissions":[
			{"service":"users","action":"read","resource":"*"},
			{"service":"orders","action":"write","resource":"*"}
		]
	}`
	req, _ := http.NewRequest("POST", ts.URL+"/v1/roles", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}
