package model

import "time"

// Role represents a named set of permissions.
type Role struct {
	ID          string
	Name        string       // unique key: admin, customer, warehouse_staff, billing_manager
	Description string
	Permissions []Permission
	IsSystem    bool // system roles cannot be deleted
	CreatedAt   time.Time
}
