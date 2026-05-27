package repository

import (
	"context"

	"github.com/beabys/wms/inventory-service/internal/domain/inventory/model"
)

// ProductRepository defines the persistence contract for Product aggregates.
type ProductRepository interface {
	// GetByID retrieves a product by ID.
	GetByID(ctx context.Context, id string) (*model.Product, error)

	// Save persists a new product.
	Save(ctx context.Context, product *model.Product) error

	// Update updates an existing product.
	Update(ctx context.Context, product *model.Product) error

	// List returns products matching the given filters.
	List(ctx context.Context, category, search string, pageSize int32, pageToken string) ([]*model.Product, string, error)

	// Archive soft-deletes a product.
	Archive(ctx context.Context, id string) error
}

// StockRepository defines the persistence contract for StockEntry aggregates.
type StockRepository interface {
	// GetByID retrieves a stock entry by ID.
	GetByID(ctx context.Context, id string) (*model.StockEntry, error)

	// Save persists a new stock entry.
	Save(ctx context.Context, entry *model.StockEntry) error

	// UpdateQuantity updates the quantity and reserved_quantity of a stock entry.
	UpdateQuantity(ctx context.Context, entry *model.StockEntry) error

	// List returns stock entries matching the given filters.
	List(ctx context.Context, productID, binLocationID, status string, pageSize int32, pageToken string) ([]*model.StockEntry, string, error)
}

// BinLocationRepository defines the persistence contract for BinLocation aggregates.
type BinLocationRepository interface {
	// GetByID retrieves a bin location by ID.
	GetByID(ctx context.Context, id string) (*model.BinLocation, error)

	// Save persists a new bin location.
	Save(ctx context.Context, location *model.BinLocation) error

	// List returns all bin locations.
	List(ctx context.Context, pageSize int32, pageToken string) ([]*model.BinLocation, string, error)
}
