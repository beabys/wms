package command

import "github.com/beabys/wms/auth-service/internal/domain/auth/model"

// CreateUserCommand represents the data needed to create a user.
type CreateUserCommand struct {
	Email      string
	Password   string
	Name       string
	Role       string
	CustomerID *string
	CreatedBy  string // admin user ID who created this user
}

// UpdateUserCommand represents the data needed to update a user.
type UpdateUserCommand struct {
	UserID string
	Name   *string
	Role   *string
	Active *bool
}

// ListUsersQuery represents the query parameters for listing users.
type ListUsersQuery struct {
	CustomerID *string
	Role       *string
	Active     *bool
	Page       int
	PageSize   int
}

// ListUsersResult contains the result of a list users query.
type ListUsersResult struct {
	Users      []UserResult
	TotalCount int
	Page       int
	PageSize   int
}

// UserResult contains user data without sensitive fields.
type UserResult struct {
	ID         string
	Email      string
	Name       string
	Role       string
	CustomerID *string
	Active     bool
	CreatedAt  string
	UpdatedAt  string
}

// AssignRoleCommand represents the data needed to assign a role.
type AssignRoleCommand struct {
	UserID string
	Role   string
}

// CreateRoleCommand represents the data needed to create a role.
type CreateRoleCommand struct {
	Name        string
	Description string
	Permissions []PermissionDTO
}

// PermissionDTO is a data transfer object for permission.
type PermissionDTO struct {
	Service  model.Service
	Action   model.Action
	Resource model.Resource
}
