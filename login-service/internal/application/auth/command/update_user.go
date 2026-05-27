package command

// UpdateUserCommand contains fields to update a user.
type UpdateUserCommand struct {
	UserID     string
	Email      string
	Role       string
	CustomerID string
}

// UpdateUserResponse contains the result of an update user operation.
type UpdateUserResponse struct {
	UserID     string
	Email      string
	Role       string
	Status     string
	CustomerID string
	UpdatedAt  string
}
