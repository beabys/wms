package command

// RefreshTokenCommand represents a token refresh request.
type RefreshTokenCommand struct {
	RefreshToken string
}

// RefreshTokenResponse contains the result of a token refresh operation.
type RefreshTokenResponse struct {
	AccessToken  string
	RefreshToken string
}
