package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCustomerStatus_String(t *testing.T) {
	tests := []struct {
		status CustomerStatus
		want   string
	}{
		{CustomerStatusPending, "pending"},
		{CustomerStatusActive, "active"},
		{CustomerStatusRejected, "rejected"},
		{CustomerStatusSuspended, "suspended"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.String())
		})
	}
}

func TestCustomerStatus_IsValid(t *testing.T) {
	assert.True(t, CustomerStatusPending.IsValid())
	assert.True(t, CustomerStatusActive.IsValid())
	assert.True(t, CustomerStatusRejected.IsValid())
	assert.True(t, CustomerStatusSuspended.IsValid())

	var bad CustomerStatus = "unknown"
	assert.False(t, bad.IsValid())
	assert.False(t, CustomerStatus("").IsValid())
}

func TestCustomerCreation(t *testing.T) {
	now := time.Now()
	c := &Customer{
		ID:             "test-id",
		CompanyName:    "Test Corp",
		Email:          "test@corp.com",
		Phone:          "1234567890",
		VatNumber:      "VAT123",
		Address:        "123 Main St",
		City:           "New York",
		PostalCode:     "10001",
		Country:        "US",
		Status:         CustomerStatusPending,
		CompanyAdminID: "admin-id",
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	assert.Equal(t, "test-id", c.ID)
	assert.Equal(t, "Test Corp", c.CompanyName)
	assert.Equal(t, "test@corp.com", c.Email)
	assert.Equal(t, "1234567890", c.Phone)
	assert.Equal(t, "VAT123", c.VatNumber)
	assert.Equal(t, "123 Main St", c.Address)
	assert.Equal(t, "New York", c.City)
	assert.Equal(t, "10001", c.PostalCode)
	assert.Equal(t, "US", c.Country)
	assert.Equal(t, CustomerStatusPending, c.Status)
	assert.Equal(t, "admin-id", c.CompanyAdminID)
	assert.Equal(t, now, c.CreatedAt)
	assert.Equal(t, now, c.UpdatedAt)
}

func TestCustomerRejectReason(t *testing.T) {
	c := &Customer{
		ID:           "test-id",
		Status:       CustomerStatusRejected,
		RejectReason: "Invalid documentation",
	}
	assert.Equal(t, "Invalid documentation", c.RejectReason)
}

func TestCompanyCreation(t *testing.T) {
	now := time.Now()
	c := &Company{
		ID:         "company-id",
		CustomerID: "customer-id",
		Name:       "Test Company",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	assert.Equal(t, "company-id", c.ID)
	assert.Equal(t, "customer-id", c.CustomerID)
	assert.Equal(t, "Test Company", c.Name)
	assert.Equal(t, now, c.CreatedAt)
	assert.Equal(t, now, c.UpdatedAt)
}

func TestCompanyRoleCreation(t *testing.T) {
	now := time.Now()
	r := &CompanyRole{
		ID:          "role-id",
		CustomerID:  "customer-id",
		Name:        "admin",
		Permissions: []string{"orders:read", "orders:write", "inventory:read"},
		IsDefault:   true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	assert.Equal(t, "role-id", r.ID)
	assert.Equal(t, "customer-id", r.CustomerID)
	assert.Equal(t, "admin", r.Name)
	assert.Equal(t, []string{"orders:read", "orders:write", "inventory:read"}, r.Permissions)
	assert.True(t, r.IsDefault)
	assert.Equal(t, now, r.CreatedAt)
	assert.Equal(t, now, r.UpdatedAt)
}

func TestCompanyRolePermissionsEmpty(t *testing.T) {
	r := &CompanyRole{
		ID:          "role-id",
		CustomerID:  "customer-id",
		Name:        "viewer",
		Permissions: []string{},
		IsDefault:   false,
	}
	assert.Empty(t, r.Permissions)
	assert.False(t, r.IsDefault)
}
