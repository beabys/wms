package command

import "github.com/beabys/wms/inventory-service/internal/domain/inventory/model"

// CreateProductCommand creates a new product.
type CreateProductCommand struct {
	SKU               string
	Name              string
	Description       string
	Category          string
	Unit              string
	WeightKg          float64
	LowStockThreshold float64
}

// ProductResult is the result of product operations.
type ProductResult struct {
	Product *model.Product
}
