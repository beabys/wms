package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomerResponseJSON(t *testing.T) {
	r := CustomerResponse{
		ID:             "cust-1",
		CompanyName:    "ACME Corp",
		Email:          "admin@acme.com",
		Phone:          "123456789",
		VatNumber:      "VAT123",
		Address:        "123 Main St",
		City:           "Springfield",
		PostalCode:     "12345",
		Country:        "US",
		Status:         "active",
		CompanyAdminID: "user-1",
		CreatedAt:      "2024-01-01T00:00:00Z",
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded CustomerResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, r, decoded)
}

func TestRegisterCustomerRequestJSON(t *testing.T) {
	r := RegisterCustomerRequest{
		Token:       "invite-token",
		CompanyName: "ACME Corp",
		Email:       "admin@acme.com",
		Password:    "secure-password",
		Phone:       "123456789",
		VatNumber:   "VAT123",
		Address:     "123 Main St",
		City:        "Springfield",
		PostalCode:  "12345",
		Country:     "US",
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded RegisterCustomerRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "invite-token", decoded.Token)
	assert.Equal(t, "ACME Corp", decoded.CompanyName)
}

func TestRegisterCustomerResponseJSON(t *testing.T) {
	r := RegisterCustomerResponse{
		Customer: &CustomerResponse{
			ID:          "cust-1",
			CompanyName: "ACME Corp",
		},
		AccessToken: "access-token-123",
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded RegisterCustomerResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "access-token-123", decoded.AccessToken)
	require.NotNil(t, decoded.Customer)
	assert.Equal(t, "cust-1", decoded.Customer.ID)
}

func TestCustomerListResponseJSON(t *testing.T) {
	r := CustomerListResponse{
		Customers: []CustomerResponse{
			{ID: "c1", CompanyName: "Company 1"},
			{ID: "c2", CompanyName: "Company 2"},
		},
		Pagination: Pagination{
			Page:       1,
			PageSize:   10,
			TotalItems: 2,
		},
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded CustomerListResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Len(t, decoded.Customers, 2)
	assert.Equal(t, 10, decoded.Pagination.PageSize)
}

func TestUpdateCustomerRequestJSON(t *testing.T) {
	r := UpdateCustomerRequest{
		Phone:      "987654321",
		Address:    "456 Oak Ave",
		City:       "Portland",
		PostalCode: "67890",
		Country:    "US",
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded UpdateCustomerRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "987654321", decoded.Phone)
}

func TestAssignCompanyRoleRequestJSON(t *testing.T) {
	r := AssignCompanyRoleRequest{
		CustomerID:  "cust-1",
		UserID:      "user-1",
		RoleName:    "admin",
		Permissions: []string{"read", "write"},
		AssignedBy:  "admin-user",
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded AssignCompanyRoleRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "admin", decoded.RoleName)
	assert.Len(t, decoded.Permissions, 2)
}

func TestCompanyRoleResponseJSON(t *testing.T) {
	r := CompanyRoleResponse{
		ID:          "role-1",
		CustomerID:  "cust-1",
		Name:        "editor",
		Permissions: []string{"read", "write"},
		IsDefault:   false,
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded CompanyRoleResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "editor", decoded.Name)
}

func TestCompanyRoleListResponseJSON(t *testing.T) {
	r := CompanyRoleListResponse{
		Roles: []CompanyRoleResponse{
			{ID: "r1", Name: "admin"},
			{ID: "r2", Name: "viewer"},
		},
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded CompanyRoleListResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Len(t, decoded.Roles, 2)
}

func TestPermissionsResponseJSON(t *testing.T) {
	r := PermissionsResponse{
		Permissions: []string{"customer:read", "customer:write"},
	}

	data, err := json.Marshal(r)
	require.NoError(t, err)

	var decoded PermissionsResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Len(t, decoded.Permissions, 2)
	assert.Equal(t, "customer:read", decoded.Permissions[0])
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
