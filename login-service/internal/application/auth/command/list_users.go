package command

// ListUsersQuery represents a request to list users with pagination.
type ListUsersQuery struct {
	Page     int32
	PageSize int32
}

// ListUsersResponse contains the result of a list users operation.
type ListUsersResponse struct {
	Users []*UserItem
}

// UserItem represents a user in a list response.
type UserItem struct {
	UserID     string
	Email      string
	Role       string
	Status     string
	CustomerID string
	CreatedAt  string
}
