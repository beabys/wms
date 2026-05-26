package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRole_Permissions(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		hasPerms bool
	}{
		{"admin", RoleAdmin, true},
		{"warehouse_staff", RoleWarehouseStaff, true},
		{"billing_manager", RoleBillingManager, true},
		{"customer", RoleCustomer, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perms := tt.role.Permissions()
			if tt.hasPerms {
				assert.NotEmpty(t, perms)
			}
		})
	}
}

func TestRole_IsValid(t *testing.T) {
	assert.True(t, RoleAdmin.IsValid())
	assert.True(t, RoleWarehouseStaff.IsValid())
	assert.True(t, RoleBillingManager.IsValid())
	assert.True(t, RoleCustomer.IsValid())
	assert.False(t, Role("invalid").IsValid())
	assert.False(t, Role("").IsValid())
}

func TestParseRole(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Role
		wantErr bool
	}{
		{"admin", "admin", RoleAdmin, false},
		{"customer", "customer", RoleCustomer, false},
		{"invalid", "superadmin", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRole(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
