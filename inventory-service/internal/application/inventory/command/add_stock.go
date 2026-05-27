package command

import (
	"time"

	"github.com/beabys/wms/inventory-service/internal/domain/inventory/model"
)

// AddStockCommand adds stock for a product.
type AddStockCommand struct {
	ProductID     string
	BinLocationID string
	Quantity      float64
	LotNumber     string
	ExpiryDate    *time.Time
}

// StockResult is the result of stock operations.
type StockResult struct {
	StockEntry *model.StockEntry
}
