package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCompanyRoleStruct(t *testing.T) {
	now := time.Now()
	r := &CompanyRole{
		ID:          "role-1",
		CustomerID:  "cust-1",
		Name:        "admin",
		Permissions: []string{"orders:read", "orders:write"},
		IsDefault:   true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	assert.Equal(t, "role-1", r.ID)
	assert.Equal(t, "cust-1", r.CustomerID)
	assert.Equal(t, "admin", r.Name)
	assert.Equal(t, []string{"orders:read", "orders:write"}, r.Permissions)
	assert.True(t, r.IsDefault)
}

func TestCompanyRoleEmptyPermissions(t *testing.T) {
	r := &CompanyRole{
		Name:        "viewer",
		Permissions: []string{},
	}
	assert.Empty(t, r.Permissions)
}
