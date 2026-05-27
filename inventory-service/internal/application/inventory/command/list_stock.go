package command

import "github.com/beabys/wms/inventory-service/internal/domain/inventory/model"

// ListStockQuery filters stock entries.
type ListStockQuery struct {
	PageSize      int32
	PageToken     string
	ProductID     string
	BinLocationID string
	Status        string
}

// ListStockResult is the result of listing stock entries.
type ListStockResult struct {
	StockEntries  []*model.StockEntry
	NextPageToken string
}
