package command

// AdjustStockCommand adjusts stock quantity.
type AdjustStockCommand struct {
	StockEntryID string
	Delta        float64
	Reason       string
	Notes        string
}
