package command

// CreateUserRequest wraps parameters for creating a user.
type CreateUserRequest struct {
	Email      string
	Password   string
	Name       string
	Role       string
	CustomerID string
	CreatedBy  string
}

// CreateUserResponse wraps the response from creating a user.
type CreateUserResponse struct {
	UserID string
	Email  string
	Name   string
	Role   string
}

// ValidateTokenResult wraps the result of token validation.
type ValidateTokenResult struct {
	UserID      string
	Role        string
	CustomerID  string
	Permissions []string
}
