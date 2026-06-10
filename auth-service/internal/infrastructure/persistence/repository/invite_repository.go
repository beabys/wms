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
	"github.com/google/uuid"
)

// inviteRow maps to the invite_tokens table in PostgreSQL.
type inviteRow struct {
	ID        string     `db:"id"`
	Email     string     `db:"email"`
	Token     string     `db:"token"`
	InvitedBy *string    `db:"invited_by"`
	Status    string     `db:"status"`
	ExpiresAt time.Time  `db:"expires_at"`
	CreatedAt time.Time  `db:"created_at"`
}

// InviteRepository implements repository.InviteRepository using PostgreSQL.
type InviteRepository struct {
	db  *sqlx.DB
	log logger.Logger
}

// NewInviteRepository creates a new InviteRepository.
func NewInviteRepository(log logger.Logger, d database.Database) *InviteRepository {
	pg, ok := d.GetDBImpl().(*database.Postgres)
	if !ok {
		return &InviteRepository{db: nil, log: log}
	}
	return &InviteRepository{db: pg.DB, log: log}
}

// Create inserts a new invite token into the database.
func (r *InviteRepository) Create(ctx context.Context, invite *model.InviteToken) error {
	if invite.ID == "" {
		invite.ID = uuid.New().String()
	}

	query := `INSERT INTO invite_tokens (id, email, token, invited_by, status, expires_at, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		invite.ID,
		invite.Email,
		invite.Token,
		invite.InvitedBy,
		invite.Status,
		invite.ExpiresAt,
		invite.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create invite token: %w", err)
	}
	return nil
}

// GetByToken retrieves an invite token by its value.
func (r *InviteRepository) GetByToken(ctx context.Context, token string) (*model.InviteToken, error) {
	var row inviteRow
	err := r.db.GetContext(ctx, &row, "SELECT * FROM invite_tokens WHERE token = $1", token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get invite token: %w", err)
	}
	return rowToInvite(&row), nil
}

// MarkUsed marks an invite token as used.
func (r *InviteRepository) MarkUsed(ctx context.Context, token string) error {
	query := "UPDATE invite_tokens SET status='used' WHERE token=$1 AND status='pending'"
	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to mark invite as used: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("invite token not found or already used/cancelled")
	}
	return nil
}

// Cancel sets invite status to cancelled (only if status=pending).
func (r *InviteRepository) Cancel(ctx context.Context, token string) error {
	query := "UPDATE invite_tokens SET status='cancelled' WHERE token=$1 AND status='pending'"
	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to cancel invite token: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("invite token not found or already processed")
	}
	return nil
}

// List retrieves paginated invite tokens with optional filters.
func (r *InviteRepository) List(ctx context.Context, filter repository.InviteFilter) (*repository.InviteListResult, error) {
	query := `SELECT id, email, token, invited_by, status, expires_at, created_at,
	           COUNT(*) OVER() as total_items
	          FROM invite_tokens WHERE 1=1`
	args := make([]interface{}, 0)
	argIdx := 1

	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, *filter.Status)
		argIdx++
	}

	if filter.Expired != nil {
		if *filter.Expired {
			query += fmt.Sprintf(" AND expires_at < NOW()")
		} else {
			query += fmt.Sprintf(" AND expires_at >= NOW()")
		}
	}

	if filter.CreatedAfter != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argIdx)
		args = append(args, *filter.CreatedAfter)
		argIdx++
	}

	if filter.CreatedBefore != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argIdx)
		args = append(args, *filter.CreatedBefore)
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	// Pagination
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list invite tokens: %w", err)
	}
	defer rows.Close()

	var invites []*model.InviteToken
	totalItems := 0
	firstRow := true

	for rows.Next() {
		var row struct {
			inviteRow
			TotalItems int `db:"total_items"`
		}
		if err := rows.StructScan(&row); err != nil {
			return nil, fmt.Errorf("failed to scan invite token row: %w", err)
		}
		if firstRow {
			totalItems = row.TotalItems
			firstRow = false
		}
		invites = append(invites, rowToInvite(&row.inviteRow))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	if firstRow {
		// No rows returned
		invites = make([]*model.InviteToken, 0)
	}

	return &repository.InviteListResult{
		Invites:    invites,
		TotalItems: totalItems,
	}, nil
}

// ExistsPendingByEmail checks if there is a pending invite for the given email.
func (r *InviteRepository) ExistsPendingByEmail(ctx context.Context, email string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM invite_tokens WHERE email=$1 AND status='pending')"
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, email)
	if err != nil {
		return false, fmt.Errorf("failed to check pending invite: %w", err)
	}
	return exists, nil
}

// Ensure InviteRepository implements repository.InviteRepository.
var _ repository.InviteRepository = (*InviteRepository)(nil)

func rowToInvite(row *inviteRow) *model.InviteToken {
	invitedBy := ""
	if row.InvitedBy != nil {
		invitedBy = *row.InvitedBy
	}
	return &model.InviteToken{
		ID:        row.ID,
		Email:     row.Email,
		Token:     row.Token,
		InvitedBy: invitedBy,
		Status:    row.Status,
		ExpiresAt: row.ExpiresAt,
		CreatedAt: row.CreatedAt,
	}
}
