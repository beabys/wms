package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/beabys/wms/inventory-service/internal/domain/inventory/model"
)

// StockRepository is the PostgreSQL implementation of application StockRepository.
type StockRepository struct {
	pool *pgxpool.Pool
}

// NewStockRepository creates a new StockRepository.
func NewStockRepository(pool *pgxpool.Pool) *StockRepository {
	return &StockRepository{pool: pool}
}

// GetByID retrieves a stock entry by ID.
func (r *StockRepository) GetByID(ctx context.Context, id string) (*model.StockEntry, error) {
	query := `
		SELECT id, product_id, bin_location_id, quantity, reserved_quantity, lot_number, expiry_date, status, last_updated
		FROM stock_entries WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)

	s := &model.StockEntry{}
	err := row.Scan(
		&s.ID, &s.ProductID, &s.BinLocationID, &s.Quantity, &s.ReservedQuantity,
		&s.LotNumber, &s.ExpiryDate, &s.Status, &s.LastUpdated,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("stock entry not found: %s", id)
		}
		return nil, fmt.Errorf("scan stock entry: %w", err)
	}
	return s, nil
}

// Save persists a new stock entry.
func (r *StockRepository) Save(ctx context.Context, s *model.StockEntry) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO stock_entries (id, product_id, bin_location_id, quantity, reserved_quantity, lot_number, expiry_date, status, last_updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, s.ID, s.ProductID, nullString(s.BinLocationID), s.Quantity, s.ReservedQuantity,
		s.LotNumber, s.ExpiryDate, s.Status, s.LastUpdated)
	if err != nil {
		return fmt.Errorf("insert stock entry: %w", err)
	}
	return nil
}

// UpdateQuantity updates the quantity and reserved_quantity.
func (r *StockRepository) UpdateQuantity(ctx context.Context, s *model.StockEntry) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE stock_entries SET quantity = $1, reserved_quantity = $2, last_updated = $3 WHERE id = $4
	`, s.Quantity, s.ReservedQuantity, s.LastUpdated, s.ID)
	if err != nil {
		return fmt.Errorf("update stock entry quantity: %w", err)
	}
	return nil
}

// List returns stock entries with optional filters.
func (r *StockRepository) List(ctx context.Context, productID, binLocationID, status string, pageSize int32, pageToken string) ([]*model.StockEntry, string, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	args := make([]any, 0, 4)
	where := "WHERE 1=1"
	argIdx := 1

	if productID != "" {
		where += fmt.Sprintf(" AND product_id = $%d", argIdx)
		args = append(args, productID)
		argIdx++
	}
	if binLocationID != "" {
		where += fmt.Sprintf(" AND bin_location_id = $%d", argIdx)
		args = append(args, binLocationID)
		argIdx++
	}
	if status != "" {
		where += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	query := fmt.Sprintf(`
		SELECT id, product_id, bin_location_id, quantity, reserved_quantity, lot_number, expiry_date, status, last_updated
		FROM stock_entries %s
		ORDER BY last_updated DESC
		LIMIT $%d
	`, where, argIdx)
	args = append(args, pageSize+1)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("list stock entries: %w", err)
	}
	defer rows.Close()

	var entries []*model.StockEntry
	for rows.Next() {
		s := &model.StockEntry{}
		err := rows.Scan(
			&s.ID, &s.ProductID, &s.BinLocationID, &s.Quantity, &s.ReservedQuantity,
			&s.LotNumber, &s.ExpiryDate, &s.Status, &s.LastUpdated,
		)
		if err != nil {
			return nil, "", fmt.Errorf("scan stock entry: %w", err)
		}
		entries = append(entries, s)
	}

	var nextToken string
	if len(entries) > int(pageSize) {
		entries = entries[:pageSize]
		nextToken = entries[len(entries)-1].ID
	}

	return entries, nextToken, nil
}

func nullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
