package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/beabys/wms/inbound-service/internal/domain/inbound/model"
)

// PGPool defines the minimal pool interface used by InboundRepository.
// Satisfied by *pgxpool.Pool. Extracted for testability.
type PGPool interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// InboundRepository is the PostgreSQL implementation of application InboundRepository.
type InboundRepository struct {
	pool PGPool
}

// NewInboundRepository creates a new InboundRepository.
func NewInboundRepository(pool *pgxpool.Pool) *InboundRepository {
	return &InboundRepository{pool: pool}
}

// GetByID retrieves an inbound with its items, inspection, and hold.
func (r *InboundRepository) GetByID(ctx context.Context, id string) (*model.Inbound, error) {
	query := `
		SELECT id, customer_id, status, expected_date, packaging_profile_id, notes, created_at, updated_at
		FROM inbounds WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)

	in := &model.Inbound{}
	var packagingProfileID *string
	err := row.Scan(
		&in.ID, &in.CustomerID, &in.Status, &in.ExpectedDate,
		&packagingProfileID, &in.Notes, &in.CreatedAt, &in.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("inbound not found: %s", id)
		}
		return nil, fmt.Errorf("scan inbound: %w", err)
	}
	if packagingProfileID != nil {
		in.PackagingProfileID = *packagingProfileID
	}

	// Load items
	items, err := r.loadItems(ctx, id)
	if err != nil {
		return nil, err
	}
	in.Items = items

	// Load inspection
	insp, err := r.loadInspection(ctx, id)
	if err != nil {
		return nil, err
	}
	in.Inspection = insp

	// Load hold
	hold, err := r.loadHold(ctx, id)
	if err != nil {
		return nil, err
	}
	in.HoldRecord = hold

	return in, nil
}

// Save persists a new inbound aggregate with items.
func (r *InboundRepository) Save(ctx context.Context, inbound *model.Inbound) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	err = r.insertInbound(ctx, tx, inbound)
	if err != nil {
		return err
	}

	for i := range inbound.Items {
		inbound.Items[i].InboundID = inbound.ID
		err = r.insertItem(ctx, tx, &inbound.Items[i])
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// UpdateStatus updates the inbound status and saves related entities.
func (r *InboundRepository) UpdateStatus(ctx context.Context, inbound *model.Inbound) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	err = r.updateInboundStatus(ctx, tx, inbound)
	if err != nil {
		return err
	}

	// Save inspection if present
	if inbound.Inspection != nil {
		err = r.upsertInspection(ctx, tx, inbound.ID, inbound.Inspection.ID, inbound.Inspection)
		if err != nil {
			return err
		}
	}

	// Save hold if present
	if inbound.HoldRecord != nil {
		err = r.upsertHold(ctx, tx, inbound.ID, inbound.HoldRecord.ID, inbound.HoldRecord)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// List returns inbounds with optional filters.
func (r *InboundRepository) List(ctx context.Context, customerID, status string, pageSize int32, pageToken string) ([]*model.Inbound, string, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	args := make([]any, 0, 4)
	where := "WHERE 1=1"
	argIdx := 1

	if customerID != "" {
		where += fmt.Sprintf(" AND customer_id = $%d", argIdx)
		args = append(args, customerID)
		argIdx++
	}
	if status != "" {
		where += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	query := fmt.Sprintf(`
		SELECT id, customer_id, status, expected_date, packaging_profile_id, notes, created_at, updated_at
		FROM inbounds %s
		ORDER BY expected_date ASC, created_at DESC
		LIMIT $%d
	`, where, argIdx)
	args = append(args, pageSize+1)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("list inbounds: %w", err)
	}
	defer rows.Close()

	var inbounds []*model.Inbound
	for rows.Next() {
		in := &model.Inbound{}
		var packagingProfileID *string
		err := rows.Scan(
			&in.ID, &in.CustomerID, &in.Status, &in.ExpectedDate,
			&packagingProfileID, &in.Notes, &in.CreatedAt, &in.UpdatedAt,
		)
		if err != nil {
			return nil, "", fmt.Errorf("scan inbound: %w", err)
		}
		if packagingProfileID != nil {
			in.PackagingProfileID = *packagingProfileID
		}
		inbounds = append(inbounds, in)
	}

	var nextToken string
	if len(inbounds) > int(pageSize) {
		inbounds = inbounds[:pageSize]
		nextToken = inbounds[len(inbounds)-1].ID
	}

	return inbounds, nextToken, nil
}

// --- internal helpers ---

func (r *InboundRepository) insertInbound(ctx context.Context, tx pgx.Tx, in *model.Inbound) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO inbounds (id, customer_id, status, expected_date, packaging_profile_id, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, in.ID, in.CustomerID, string(in.Status), in.ExpectedDate,
		nullString(in.PackagingProfileID), in.Notes, in.CreatedAt, in.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert inbound: %w", err)
	}
	return nil
}

func (r *InboundRepository) insertItem(ctx context.Context, tx pgx.Tx, item *model.InboundItem) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO inbound_items (id, inbound_id, sku, quantity_declared, quantity_received, dimensions, weight)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, item.ID, item.InboundID, item.SKU, item.QuantityDeclared, item.QuantityReceived,
		item.Dimensions, item.Weight)
	if err != nil {
		return fmt.Errorf("insert item: %w", err)
	}
	return nil
}

func (r *InboundRepository) updateInboundStatus(ctx context.Context, tx pgx.Tx, in *model.Inbound) error {
	_, err := tx.Exec(ctx, `
		UPDATE inbounds SET status = $1, updated_at = $2 WHERE id = $3
	`, string(in.Status), in.UpdatedAt, in.ID)
	if err != nil {
		return fmt.Errorf("update inbound: %w", err)
	}
	return nil
}

func (r *InboundRepository) upsertInspection(ctx context.Context, tx pgx.Tx, inboundID, inspectionID string, insp *model.Inspection) error {
	photosJSON, _ := json.Marshal(insp.Photos)
	_, err := tx.Exec(ctx, `
		INSERT INTO inspections (id, inbound_id, inspector_id, notes, photos_json, passed, inspected_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (inbound_id) DO UPDATE SET
			inspector_id = EXCLUDED.inspector_id,
			notes = EXCLUDED.notes,
			photos_json = EXCLUDED.photos_json,
			passed = EXCLUDED.passed,
			inspected_at = EXCLUDED.inspected_at
	`, inspectionID, inboundID, insp.InspectorID, insp.Notes, string(photosJSON), insp.Passed, insp.InspectedAt)
	if err != nil {
		return fmt.Errorf("upsert inspection: %w", err)
	}
	return nil
}

func (r *InboundRepository) upsertHold(ctx context.Context, tx pgx.Tx, inboundID, holdID string, hold *model.Hold) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO holds (id, inbound_id, reason, created_at, released_at, released_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (inbound_id) DO UPDATE SET
			reason = EXCLUDED.reason,
			created_at = EXCLUDED.created_at,
			released_at = EXCLUDED.released_at,
			released_by = EXCLUDED.released_by
	`, holdID, inboundID, hold.Reason, hold.CreatedAt, hold.ReleasedAt, hold.ReleasedBy)
	if err != nil {
		return fmt.Errorf("upsert hold: %w", err)
	}
	return nil
}

func (r *InboundRepository) loadItems(ctx context.Context, inboundID string) ([]model.InboundItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, inbound_id, sku, quantity_declared, quantity_received, dimensions, weight
		FROM inbound_items WHERE inbound_id = $1
	`, inboundID)
	if err != nil {
		return nil, fmt.Errorf("load items: %w", err)
	}
	defer rows.Close()

	var items []model.InboundItem
	for rows.Next() {
		var item model.InboundItem
		err := rows.Scan(&item.ID, &item.InboundID, &item.SKU,
			&item.QuantityDeclared, &item.QuantityReceived,
			&item.Dimensions, &item.Weight)
		if err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *InboundRepository) loadInspection(ctx context.Context, inboundID string) (*model.Inspection, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, inspector_id, notes, photos_json, passed, inspected_at
		FROM inspections WHERE inbound_id = $1
	`, inboundID)

	var insp model.Inspection
	var photosJSON string
	err := row.Scan(&insp.ID, &insp.InspectorID, &insp.Notes, &photosJSON, &insp.Passed, &insp.InspectedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan inspection: %w", err)
	}
	_ = json.Unmarshal([]byte(photosJSON), &insp.Photos)
	return &insp, nil
}

func (r *InboundRepository) loadHold(ctx context.Context, inboundID string) (*model.Hold, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, reason, created_at, released_at, released_by
		FROM holds WHERE inbound_id = $1
	`, inboundID)

	var hold model.Hold
	err := row.Scan(&hold.ID, &hold.Reason, &hold.CreatedAt, &hold.ReleasedAt, &hold.ReleasedBy)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan hold: %w", err)
	}
	return &hold, nil
}

func nullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

