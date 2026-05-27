package command

// CreateUserCommand represents a request to create a new user by an admin.
type CreateUserCommand struct {
	Email      string
	Password   string
	Role       string
	CustomerID string
}

// CreateUserResponse contains the result of a create user operation.
type CreateUserResponse struct {
	UserID    string
	Email     string
	Role      string
	Status    string
	CreatedAt string
}
