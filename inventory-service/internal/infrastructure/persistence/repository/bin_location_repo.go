package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/beabys/wms/inventory-service/internal/domain/inventory/model"
)

// BinLocationRepository is the PostgreSQL implementation of application BinLocationRepository.
type BinLocationRepository struct {
	pool *pgxpool.Pool
}

// NewBinLocationRepository creates a new BinLocationRepository.
func NewBinLocationRepository(pool *pgxpool.Pool) *BinLocationRepository {
	return &BinLocationRepository{pool: pool}
}

// GetByID retrieves a bin location by ID.
func (r *BinLocationRepository) GetByID(ctx context.Context, id string) (*model.BinLocation, error) {
	query := `
		SELECT id, warehouse_zone, aisle, rack, shelf, is_active
		FROM bin_locations WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)

	b := &model.BinLocation{}
	err := row.Scan(&b.ID, &b.WarehouseZone, &b.Aisle, &b.Rack, &b.Shelf, &b.IsActive)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("bin location not found: %s", id)
		}
		return nil, fmt.Errorf("scan bin location: %w", err)
	}
	return b, nil
}

// Save persists a new bin location.
func (r *BinLocationRepository) Save(ctx context.Context, b *model.BinLocation) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO bin_locations (id, warehouse_zone, aisle, rack, shelf, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, b.ID, b.WarehouseZone, b.Aisle, b.Rack, b.Shelf, b.IsActive)
	if err != nil {
		return fmt.Errorf("insert bin location: %w", err)
	}
	return nil
}

// List returns all bin locations.
func (r *BinLocationRepository) List(ctx context.Context, pageSize int32, pageToken string) ([]*model.BinLocation, string, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	query := `
		SELECT id, warehouse_zone, aisle, rack, shelf, is_active
		FROM bin_locations
		ORDER BY warehouse_zone, aisle, rack, shelf
		LIMIT $1
	`
	rows, err := r.pool.Query(ctx, query, pageSize+1)
	if err != nil {
		return nil, "", fmt.Errorf("list bin locations: %w", err)
	}
	defer rows.Close()

	var locations []*model.BinLocation
	for rows.Next() {
		b := &model.BinLocation{}
		err := rows.Scan(&b.ID, &b.WarehouseZone, &b.Aisle, &b.Rack, &b.Shelf, &b.IsActive)
		if err != nil {
			return nil, "", fmt.Errorf("scan bin location: %w", err)
		}
		locations = append(locations, b)
	}

	var nextToken string
	if len(locations) > int(pageSize) {
		locations = locations[:pageSize]
		nextToken = locations[len(locations)-1].ID
	}

	return locations, nextToken, nil
}
