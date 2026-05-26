package model

import "fmt"

// Role represents a user role in the system.
type Role string

const (
	RoleAdmin          Role = "admin"
	RoleWarehouseStaff Role = "warehouse_staff"
	RoleBillingManager Role = "billing_manager"
	RoleCustomer       Role = "customer"
)

// Permissions returns the permissions associated with this role.
func (r Role) Permissions() []string {
	switch r {
	case RoleAdmin:
		return []string{"users:read", "users:write", "users:delete", "warehouse:read", "warehouse:write", "billing:read", "billing:write"}
	case RoleWarehouseStaff:
		return []string{"warehouse:read", "warehouse:write"}
	case RoleBillingManager:
		return []string{"billing:read", "billing:write", "users:read"}
	case RoleCustomer:
		return []string{"warehouse:read"}
	default:
		return nil
	}
}

// IsValid checks if the role is a valid enum value.
func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleWarehouseStaff, RoleBillingManager, RoleCustomer:
		return true
	default:
		return false
	}
}

// ParseRole parses a string into a Role.
func ParseRole(s string) (Role, error) {
	r := Role(s)
	if !r.IsValid() {
		return "", fmt.Errorf("invalid role: %s", s)
	}
	return r, nil
}
