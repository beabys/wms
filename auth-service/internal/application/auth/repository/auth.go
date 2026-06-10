package repository

import (
	"context"
	"time"
)

// AuthRepository handles refresh token storage and validation.
type AuthRepository interface {
	// StoreRefreshToken stores a refresh token hash for a user.
	StoreRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error

	// ValidateRefreshToken checks if a refresh token hash is valid and returns the userID.
	// Returns empty string if token is invalid or expired.
	ValidateRefreshToken(ctx context.Context, tokenHash string) (string, error)

	// RevokeRefreshToken marks a refresh token as revoked.
	RevokeRefreshToken(ctx context.Context, tokenHash string) error

	// RevokeAllUserTokens revokes all refresh tokens for a user (e.g., password change).
	RevokeAllUserTokens(ctx context.Context, userID string) error
}
