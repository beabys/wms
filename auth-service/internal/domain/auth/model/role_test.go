package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermissionString(t *testing.T) {
	t.Run("returns canonical format", func(t *testing.T) {
		p := Permission{Service: ServiceAuth, Action: ActionRead, Resource: ResourceUsers}
		assert.Equal(t, "auth:read:users", p.String())
	})

	t.Run("handles empty fields", func(t *testing.T) {
		p := Permission{}
		assert.Equal(t, "::", p.String())
	})
}

func TestRolePermissions(t *testing.T) {
	t.Run("role can hold multiple permissions", func(t *testing.T) {
		role := Role{
			ID:   "role-1",
			Name: "admin",
			Permissions: []Permission{
				{Service: ServiceAuth, Action: ActionAdmin, Resource: ResourceUsers},
				{Service: ServiceCustomer, Action: ActionWrite, Resource: ResourceCustomers},
			},
		}
		assert.Len(t, role.Permissions, 2)
		assert.Equal(t, "auth:admin:users", role.Permissions[0].String())
	})

	t.Run("system role flag works", func(t *testing.T) {
		sysRole := Role{ID: "sys-1", Name: "admin", IsSystem: true}
		userRole := Role{ID: "usr-1", Name: "custom", IsSystem: false}
		assert.True(t, sysRole.IsSystem)
		assert.False(t, userRole.IsSystem)
	})
}
