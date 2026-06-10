package model

// InviteUserResponse represents the response for inviting a user.
type InviteUserResponse struct {
	Token      string `json:"token"`
	InviteLink string `json:"invite_link"`
	ExpiresAt  int64  `json:"expires_at"`
}

// AuthResponse represents the authentication token response.
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// UserResponse represents a user in API responses.
type UserResponse struct {
	ID         string  `json:"id"`
	Email      string  `json:"email"`
	Name       string  `json:"name"`
	Role       string  `json:"role"`
	CustomerID *string `json:"customer_id,omitempty"`
	Active     bool    `json:"active"`
	CreatedAt  string  `json:"created_at"`
}

// UserListResponse represents a paginated list of users.
type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Pagination Pagination     `json:"pagination"`
}

// InviteEntry represents an invite token in API responses.
type InviteEntry struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	InvitedBy string `json:"invited_by"`
	Status    string `json:"status"`
	ExpiresAt int64  `json:"expires_at"`
	CreatedAt int64  `json:"created_at"`
}

// InviteListResponse represents a paginated list of invites.
type InviteListResponse struct {
	Invites    []InviteEntry `json:"invites"`
	Pagination Pagination    `json:"pagination"`
}

// CancelInviteResponse represents the response for cancelling an invite.
type CancelInviteResponse struct {
	Success bool `json:"success"`
}

// Pagination contains pagination metadata.
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
}

// CreateUserRequest represents the request body for creating a user.
type CreateUserRequest struct {
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	Name       string  `json:"name"`
	Role       string  `json:"role"`
	CustomerID *string `json:"customer_id,omitempty"`
}

// UpdateUserRequest represents the request body for updating a user.
type UpdateUserRequest struct {
	Name   *string `json:"name,omitempty"`
	Role   *string `json:"role,omitempty"`
	Active *bool   `json:"active,omitempty"`
}

// RoleResponse represents a role in API responses.
type RoleResponse struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Permissions []RolePermission `json:"permissions,omitempty"`
	IsSystem    bool             `json:"is_system"`
}

// RolePermission represents a permission entry in a role.
type RolePermission struct {
	Service  string `json:"service"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
}

// RoleListResponse represents a list of roles.
type RoleListResponse struct {
	Roles []RoleResponse `json:"roles"`
}

// CreateRoleRequest represents the request body for creating a role.
type CreateRoleRequest struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Permissions []RolePermission `json:"permissions,omitempty"`
}
