package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/beabys/wms/auth-service/internal/application/auth/repository"
	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
	"github.com/beabys/wms/pkg/database"
	"github.com/beabys/wms/pkg/logger"
	"github.com/jmoiron/sqlx"
)

// userRow maps to the users table in PostgreSQL.
type userRow struct {
	ID           string     `db:"id"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	Name         string     `db:"name"`
	Role         string     `db:"role"`
	CustomerID   *string    `db:"customer_id"`
	Active       bool       `db:"active"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

// UserRepository implements repository.UserRepository using PostgreSQL.
type UserRepository struct {
	db  *sqlx.DB
	log logger.Logger
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(log logger.Logger, d database.Database) *UserRepository {
	pg, ok := d.GetDBImpl().(*database.Postgres)
	if !ok {
		return &UserRepository{db: nil, log: log}
	}
	return &UserRepository{db: pg.DB, log: log}
}

// GetByID retrieves a user by their ID.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var row userRow
	err := r.db.GetContext(ctx, &row, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return rowToUser(&row), nil
}

// GetByEmail retrieves a user by their email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var row userRow
	err := r.db.GetContext(ctx, &row, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return rowToUser(&row), nil
}

// Create inserts a new user into the database.
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (id, email, password_hash, name, role, customer_id, active, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		string(user.Password),
		user.Name,
		user.Role,
		user.CustomerID,
		user.Active,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// Update updates an existing user.
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `UPDATE users SET email=$1, name=$2, role=$3, customer_id=$4, active=$5, updated_at=$6 WHERE id=$7`

	result, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.Name,
		user.Role,
		user.CustomerID,
		user.Active,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// List retrieves users with filtering and pagination.
func (r *UserRepository) List(ctx context.Context, filter repository.UserFilter) ([]*model.User, int, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	// Build WHERE clause
	where := "WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if filter.CustomerID != nil {
		where += fmt.Sprintf(" AND customer_id = $%d", argIdx)
		args = append(args, *filter.CustomerID)
		argIdx++
	}
	if filter.Role != nil {
		where += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, *filter.Role)
		argIdx++
	}
	if filter.Active != nil {
		where += fmt.Sprintf(" AND active = $%d", argIdx)
		args = append(args, *filter.Active)
		argIdx++
	}

	// Count
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", where)
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Fetch with pagination
	offset := (filter.Page - 1) * filter.PageSize
	dataQuery := fmt.Sprintf("SELECT * FROM users %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d", where, argIdx, argIdx+1)
	args = append(args, filter.PageSize, offset)

	var rows []userRow
	if err := r.db.SelectContext(ctx, &rows, dataQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*model.User, 0, len(rows))
	for i := range rows {
		users = append(users, rowToUser(&rows[i]))
	}

	return users, total, nil
}

// Deactivate sets a user's active flag to false.
func (r *UserRepository) Deactivate(ctx context.Context, id string) error {
	query := "UPDATE users SET active=false, updated_at=$1 WHERE id=$2"
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func rowToUser(row *userRow) *model.User {
	return &model.User{
		ID:         row.ID,
		Email:      row.Email,
		Password:   model.PasswordHash(row.PasswordHash),
		Name:       row.Name,
		Role:       row.Role,
		CustomerID: row.CustomerID,
		Active:     row.Active,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
}
