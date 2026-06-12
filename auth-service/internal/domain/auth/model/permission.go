package model

import "fmt"

// Service represents a service domain in the system.
type Service string

const (
	ServiceAuth      Service = "auth"
	ServiceCustomer  Service = "customer"
	ServiceWarehouse Service = "warehouse"
)

// Action represents an action that can be performed on a resource.
type Action string

const (
	ActionRead  Action = "read"
	ActionWrite Action = "write"
	ActionAdmin Action = "admin"
)

// Resource represents a resource type in the system.
type Resource string

const (
	ResourceUsers       Resource = "users"
	ResourceRoles       Resource = "roles"
	ResourcePermissions Resource = "permissions"
	ResourceCustomers   Resource = "customers"
)

// Permission represents a specific action on a resource within a service.
type Permission struct {
	Service  Service
	Action   Action
	Resource Resource
}

// String returns the canonical string representation: service:action:resource.
func (p Permission) String() string {
	return fmt.Sprintf("%s:%s:%s", p.Service, p.Action, p.Resource)
}
