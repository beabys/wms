package command

// LoginCommand represents a login request with email and password.
type LoginCommand struct {
	Email    string
	Password string
}

// LoginResponse contains the result of a login operation.
type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	UserID       string
	UserEmail    string
	UserRole     string
	CustomerID   string
}
