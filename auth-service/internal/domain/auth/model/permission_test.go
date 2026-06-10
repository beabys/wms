package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceConstants(t *testing.T) {
	assert.Equal(t, Service("auth"), ServiceAuth)
	assert.Equal(t, Service("customer"), ServiceCustomer)
	assert.Equal(t, Service("warehouse"), ServiceWarehouse)
}

func TestActionConstants(t *testing.T) {
	assert.Equal(t, Action("read"), ActionRead)
	assert.Equal(t, Action("write"), ActionWrite)
	assert.Equal(t, Action("admin"), ActionAdmin)
}

func TestResourceConstants(t *testing.T) {
	assert.Equal(t, Resource("users"), ResourceUsers)
	assert.Equal(t, Resource("roles"), ResourceRoles)
	assert.Equal(t, Resource("permissions"), ResourcePermissions)
	assert.Equal(t, Resource("customers"), ResourceCustomers)
}

func TestPermissionCanonicalFormat(t *testing.T) {
	tests := []struct {
		name       string
		permission Permission
		want       string
	}{
		{
			name:       "full permission with typed constants",
			permission: Permission{Service: ServiceAuth, Action: ActionRead, Resource: ResourceUsers},
			want:       "auth:read:users",
		},
		{
			name:       "admin action on customer resource",
			permission: Permission{Service: ServiceCustomer, Action: ActionAdmin, Resource: ResourceCustomers},
			want:       "customer:admin:customers",
		},
		{
			name:       "warehouse write permissions",
			permission: Permission{Service: ServiceWarehouse, Action: ActionWrite, Resource: ResourcePermissions},
			want:       "warehouse:write:permissions",
		},
		{
			name:       "empty fields",
			permission: Permission{},
			want:       "::",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.permission.String())
		})
	}
}

func TestPermissionEquality(t *testing.T) {
	p1 := Permission{Service: ServiceAuth, Action: ActionRead, Resource: ResourceUsers}
	p2 := Permission{Service: ServiceAuth, Action: ActionRead, Resource: ResourceUsers}
	p3 := Permission{Service: ServiceAuth, Action: ActionWrite, Resource: ResourceUsers}

	assert.Equal(t, p1, p2)
	assert.NotEqual(t, p1, p3)
}

func TestPermissionCreation(t *testing.T) {
	p := Permission{
		Service:  ServiceCustomer,
		Action:   ActionWrite,
		Resource: ResourceCustomers,
	}
	assert.Equal(t, ServiceCustomer, p.Service)
	assert.Equal(t, ActionWrite, p.Action)
	assert.Equal(t, ResourceCustomers, p.Resource)
}

// TestPermissionTypeSafety ensures that only valid constant values compile
// and that zero values are not valid enum values.
func TestPermissionTypeSafety(t *testing.T) {
	// Zero values should be empty strings (not valid enum values)
	var s Service
	var a Action
	var r Resource
	assert.Empty(t, string(s))
	assert.Empty(t, string(a))
	assert.Empty(t, string(r))

	// Casting from string is explicit type conversion, not validation
	// This test documents that the type system prevents accidental mixups
	// but does NOT validate at runtime — validation should happen at boundaries.
	s = Service("invalid")
	a = Action("nope")
	r = Resource("bad")
	assert.Equal(t, Service("invalid"), s)
	assert.Equal(t, Action("nope"), a)
	assert.Equal(t, Resource("bad"), r)
}
