package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/beabys/wms/pkg/database"
	"github.com/beabys/wms/pkg/logger"
	"github.com/jmoiron/sqlx"
)

// refreshTokenRow maps to the refresh_tokens table.
type refreshTokenRow struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	ExpiresAt time.Time `db:"expires_at"`
	Revoked   bool      `db:"revoked"`
	CreatedAt time.Time `db:"created_at"`
}

// AuthRepository implements repository.AuthRepository using PostgreSQL.
type AuthRepository struct {
	db  *sqlx.DB
	log logger.Logger
}

// NewAuthRepository creates a new AuthRepository.
func NewAuthRepository(log logger.Logger, d database.Database) *AuthRepository {
	pg, ok := d.GetDBImpl().(*database.Postgres)
	if !ok {
		return &AuthRepository{db: nil, log: log}
	}
	return &AuthRepository{db: pg.DB, log: log}
}

// StoreRefreshToken inserts a new refresh token record.
func (r *AuthRepository) StoreRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at, revoked, created_at)
	          VALUES ($1, $2, $3, false, $4)`
	_, err := r.db.ExecContext(ctx, query, userID, tokenHash, expiresAt, time.Now())
	if err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}
	return nil
}

// ValidateRefreshToken checks if a refresh token hash is valid and returns the userID.
func (r *AuthRepository) ValidateRefreshToken(ctx context.Context, tokenHash string) (string, error) {
	var row refreshTokenRow
	err := r.db.GetContext(ctx, &row,
		"SELECT * FROM refresh_tokens WHERE token_hash = $1 AND revoked = false AND expires_at > $2 LIMIT 1",
		tokenHash, time.Now(),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to validate refresh token: %w", err)
	}
	return row.UserID, nil
}

// RevokeRefreshToken marks a refresh token as revoked.
func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE refresh_tokens SET revoked = true WHERE token_hash = $1", tokenHash)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	return nil
}

// RevokeAllUserTokens revokes all refresh tokens for a user.
func (r *AuthRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE refresh_tokens SET revoked = true WHERE user_id = $1 AND revoked = false", userID)
	if err != nil {
		return fmt.Errorf("failed to revoke all user tokens: %w", err)
	}
	return nil
}
