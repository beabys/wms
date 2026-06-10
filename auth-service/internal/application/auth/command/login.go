package command

// LoginCommand represents the login request data.
type LoginCommand struct {
	Email    string
	Password string
}

// LoginResult contains the result of a successful login.
type LoginResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // seconds until access token expires
}

// RefreshTokenCommand represents the refresh token request data.
type RefreshTokenCommand struct {
	RefreshToken string
}

// ValidateTokenCommand represents the validate token request data.
type ValidateTokenCommand struct {
	Token string
}

// ValidateTokenResult contains claims from a validated token.
type ValidateTokenResult struct {
	UserID      string
	Role        string
	CustomerID  string
	Permissions []string
}

// LogoutCommand represents the logout request data.
type LogoutCommand struct {
	RefreshToken string
}

// GetPublicKeyResult contains the RSA public key PEM.
type GetPublicKeyResult struct {
	PublicKeyPEM string
}
