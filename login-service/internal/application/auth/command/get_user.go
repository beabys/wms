package command

// GetUserQuery represents a request to get a user by ID.
type GetUserQuery struct {
	UserID string
}

// GetUserResponse contains the result of a get user operation.
type GetUserResponse struct {
	UserID     string
	Email      string
	Role       string
	Status     string
	CustomerID string
	CreatedAt  string
	UpdatedAt  string
}
