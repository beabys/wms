package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
)

// InviteRepo implements repository.InviteLinkRepository using pgx.
type InviteRepo struct {
	pool *pgxpool.Pool
}

// NewInviteRepo creates a new PostgreSQL invite link repository.
func NewInviteRepo(pool *pgxpool.Pool) *InviteRepo {
	return &InviteRepo{pool: pool}
}

// Create inserts a new invite link.
func (r *InviteRepo) Create(ctx context.Context, link *model.InviteLink) error {
	query := `INSERT INTO invite_links (id, email, token, expires_at, used, created_at)
	           VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.pool.Exec(ctx, query,
		link.ID, link.Email, link.Token, link.ExpiresAt, link.Used, link.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("pg insert invite: %w", err)
	}
	return nil
}

// GetByToken retrieves an invite link by its token.
func (r *InviteRepo) GetByToken(ctx context.Context, token string) (*model.InviteLink, error) {
	query := `SELECT id, email, token, expires_at, used, created_at
	           FROM invite_links WHERE token = $1`

	row := r.pool.QueryRow(ctx, query, token)
	return scanInvite(row)
}

// MarkUsed sets the invite as used.
func (r *InviteRepo) MarkUsed(ctx context.Context, id string) error {
	query := `UPDATE invite_links SET used = true WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("pg mark invite used: %w", err)
	}
	return nil
}

func scanInvite(row rowScanner) (*model.InviteLink, error) {
	var (
		link      model.InviteLink
		expiresAt time.Time
		createdAt time.Time
	)

	if err := row.Scan(&link.ID, &link.Email, &link.Token, &expiresAt, &link.Used, &createdAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invite not found")
		}
		return nil, fmt.Errorf("pg scan invite: %w", err)
	}
	link.ExpiresAt = expiresAt
	link.CreatedAt = createdAt

	return &link, nil
}
