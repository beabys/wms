package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthResponseJSON(t *testing.T) {
	r := AuthResponse{
		AccessToken:  "access-token-abc",
		RefreshToken: "refresh-token-xyz",
		ExpiresIn:    3600,
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded AuthResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, r, decoded)
}

func TestUserResponseJSON(t *testing.T) {
	cid := "cust-1"
	r := UserResponse{
		ID:         "user-1",
		Email:      "user@example.com",
		Name:       "John",
		Role:       "admin",
		CustomerID: &cid,
		Active:     true,
		CreatedAt:  "2024-01-01T00:00:00Z",
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded UserResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, r, decoded)
	assert.Equal(t, "cust-1", *decoded.CustomerID)
}

func TestUserResponseNilCustomerID(t *testing.T) {
	r := UserResponse{
		ID:    "user-2",
		Email: "user2@example.com",
		Name:  "Jane",
		Role:  "viewer",
		Active: false,
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	// CustomerID should be omitted in JSON
	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	require.NoError(t, err)
	_, exists := raw["customer_id"]
	assert.False(t, exists, "customer_id should be omitted when nil")
}

func TestUserListResponseJSON(t *testing.T) {
	r := UserListResponse{
		Users: []UserResponse{
			{ID: "u1", Email: "a@b.com"},
			{ID: "u2", Email: "c@d.com"},
		},
		Pagination: Pagination{
			Page:       1,
			PageSize:   10,
			TotalItems: 2,
		},
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded UserListResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Len(t, decoded.Users, 2)
	assert.Equal(t, 10, decoded.Pagination.PageSize)
}

func TestCreateUserRequestJSON(t *testing.T) {
	cid := "cust-1"
	r := CreateUserRequest{
		Email:      "new@example.com",
		Password:   "secure-password",
		Name:       "New User",
		Role:       "admin",
		CustomerID: &cid,
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded CreateUserRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "new@example.com", decoded.Email)
	assert.Equal(t, "cust-1", *decoded.CustomerID)
}

func TestUpdateUserRequestJSON(t *testing.T) {
	name := "Updated Name"
	active := true
	r := UpdateUserRequest{
		Name:   &name,
		Active: &active,
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded UpdateUserRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", *decoded.Name)
	assert.True(t, *decoded.Active)
	assert.Nil(t, decoded.Role)
}

func TestRoleResponseJSON(t *testing.T) {
	r := RoleResponse{
		ID:          "role-1",
		Name:        "admin",
		Description: "Full access",
		Permissions: []RolePermission{
			{Service: "login", Action: "read", Resource: "users"},
		},
		IsSystem: true,
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded RoleResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Len(t, decoded.Permissions, 1)
	assert.True(t, decoded.IsSystem)
}

func TestRoleListResponseJSON(t *testing.T) {
	r := RoleListResponse{
		Roles: []RoleResponse{
			{ID: "r1", Name: "admin"},
			{ID: "r2", Name: "viewer"},
		},
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded RoleListResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Len(t, decoded.Roles, 2)
}

func TestCreateRoleRequestJSON(t *testing.T) {
	r := CreateRoleRequest{
		Name:        "editor",
		Description: "Can edit content",
		Permissions: []RolePermission{
			{Service: "cms", Action: "write", Resource: "articles"},
		},
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded CreateRoleRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "editor", decoded.Name)
	assert.Len(t, decoded.Permissions, 1)
}

func TestPaginationJSON(t *testing.T) {
	p := Pagination{
		Page:       3,
		PageSize:   25,
		TotalItems: 100,
	}

	data, err := json.Marshal(p)
	require.NoError(t, err)

	var decoded Pagination
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, 3, decoded.Page)
	assert.Equal(t, 100, decoded.TotalItems)
}
