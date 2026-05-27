package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/beabys/wms/login-service/internal/domain/auth/model"
)

// UserRepo is the pgx implementation of UserRepository.
type UserRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

// GetByID retrieves a user by their ID.
func (r *UserRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT id, email, password_hash, role, customer_id, status, created_at, updated_at FROM users WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)
	return scanUser(row)
}

// GetByEmail retrieves a user by their email.
func (r *UserRepo) GetByEmail(ctx context.Context, email model.Email) (*model.User, error) {
	query := `SELECT id, email, password_hash, role, customer_id, status, created_at, updated_at FROM users WHERE email = $1`
	row := r.pool.QueryRow(ctx, query, email.String())
	return scanUser(row)
}

// Save inserts a new user.
func (r *UserRepo) Save(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (id, email, password_hash, role, customer_id, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	_, err := r.pool.Exec(ctx, query, user.ID, user.Email.String(), user.PasswordHash.String(),
		string(user.Role), user.CustomerID, string(user.Status), user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

// List retrieves users with pagination.
func (r *UserRepo) List(ctx context.Context, limit, offset int32) ([]*model.User, error) {
	query := `SELECT id, email, password_hash, role, customer_id, status, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user, err := scanUserFromRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan user row: %w", err)
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// Update updates an existing user.
func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	query := `UPDATE users SET email = $2, role = $3, customer_id = $4, status = $5, updated_at = $6 WHERE id = $1`
	user.UpdatedAt = time.Now()
	_, err := r.pool.Exec(ctx, query, user.ID, user.Email.String(), string(user.Role),
		user.CustomerID, string(user.Status), user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// CreateRefreshToken stores a refresh token for a user.
func (r *UserRepo) CreateRefreshToken(ctx context.Context, userID, token string) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	// Parse expiry from token claims - store with 7d default
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	_, err := r.pool.Exec(ctx, query, userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("insert refresh token: %w", err)
	}
	return nil
}

// ValidateRefreshToken checks if a refresh token exists and returns the user ID.
func (r *UserRepo) ValidateRefreshToken(ctx context.Context, token string) (string, error) {
	query := `SELECT user_id FROM refresh_tokens WHERE token_hash = $1 AND expires_at > NOW()`
	var userID string
	err := r.pool.QueryRow(ctx, query, token).Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("refresh token not found or expired")
		}
		return "", fmt.Errorf("validate refresh token: %w", err)
	}
	return userID, nil
}

// DeleteRefreshToken removes a refresh token.
func (r *UserRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := r.pool.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}
	return nil
}

// scanUser scans a single user row.
func scanUser(row pgx.Row) (*model.User, error) {
	user := &model.User{}
	var emailStr, roleStr, statusStr string
	err := row.Scan(&user.ID, &emailStr, &user.PasswordHash, &roleStr, &user.CustomerID, &statusStr, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("scan user: %w", err)
	}
	user.Email = model.Email(emailStr)
	user.Role = model.Role(roleStr)
	user.Status = model.UserStatus(statusStr)
	return user, nil
}

// scannableRow is an interface satisfied by both pgx.Row and pgx.Rows.
type scannableRow interface {
	Scan(dest ...interface{}) error
}

// scanUserFromRow scans a user from a pgx.Rows iterator.
func scanUserFromRow(rows pgx.Rows) (*model.User, error) {
	user := &model.User{}
	var emailStr, roleStr, statusStr string
	err := rows.Scan(&user.ID, &emailStr, &user.PasswordHash, &roleStr, &user.CustomerID, &statusStr, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	user.Email = model.Email(emailStr)
	user.Role = model.Role(roleStr)
	user.Status = model.UserStatus(statusStr)
	return user, nil
}
