package command

// ReserveStockCommand reserves quantity from a stock entry.
type ReserveStockCommand struct {
	StockEntryID string
	Quantity     float64
}
