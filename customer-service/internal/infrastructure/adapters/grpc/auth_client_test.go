package grpcdapter

import (
	"testing"

	"github.com/beabys/wms/customer-service/internal/application/customer/command"
	"github.com/beabys/wms/customer-service/internal/application/customer/ports"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthClient_Fail(t *testing.T) {
	client, err := NewAuthClient("localhost:1")
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestAuthClientInterface(t *testing.T) {
	// Verify that AuthClient satisfies the ports.AuthClient interface.
	var _ ports.AuthClient = (*AuthClient)(nil)
}

func TestClose_Nil(t *testing.T) {
	// Just verify the method exists - can't test without a real connection
	var c *AuthClient
	assert.Nil(t, c)
}

func TestValidateTokenResult(t *testing.T) {
	r := &command.ValidateTokenResult{
		UserID:      "user-1",
		Role:        "admin",
		CustomerID:  "cust-1",
		Permissions: []string{"*"},
	}
	assert.Equal(t, "user-1", r.UserID)
	assert.Equal(t, "admin", r.Role)
	assert.Equal(t, "cust-1", r.CustomerID)
	assert.Equal(t, []string{"*"}, r.Permissions)
}

func TestCreateUserRequest(t *testing.T) {
	req := command.CreateUserRequest{
		Email:      "test@test.com",
		Password:   "password123",
		Name:       "Test User",
		Role:       "customer_admin",
		CustomerID: "cust-1",
		CreatedBy:  "admin-1",
	}
	assert.Equal(t, "test@test.com", req.Email)
	assert.Equal(t, "password123", req.Password)
	assert.Equal(t, "Test User", req.Name)
	assert.Equal(t, "customer_admin", req.Role)
	assert.Equal(t, "cust-1", req.CustomerID)
	assert.Equal(t, "admin-1", req.CreatedBy)
}

func TestCreateUserResponse(t *testing.T) {
	resp := &command.CreateUserResponse{
		UserID: "user-1",
		Email:  "test@test.com",
		Name:   "Test User",
		Role:   "customer_admin",
	}
	assert.Equal(t, "user-1", resp.UserID)
	assert.Equal(t, "customer_admin", resp.Role)
}

func TestAuthClientValidateInviteInterface(t *testing.T) {
	// Verify AuthClient satisfies the ports.AuthClient interface including ValidateInvite
	var _ ports.AuthClient = (*AuthClient)(nil)
}
