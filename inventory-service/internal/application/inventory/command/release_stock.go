package command

// ReleaseStockCommand releases reserved stock.
type ReleaseStockCommand struct {
	StockEntryID string
	Quantity     float64
}
