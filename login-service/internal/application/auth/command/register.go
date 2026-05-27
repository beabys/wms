package command

// RegisterCommand represents a user registration request.
type RegisterCommand struct {
	Email    string
	Password string
}

// RegisterResponse contains the result of a registration operation.
type RegisterResponse struct {
	UserID    string
	Email     string
	Role      string
	Status    string
	CreatedAt string
}
