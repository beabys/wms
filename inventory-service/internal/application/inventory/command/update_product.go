package command

// UpdateProductCommand updates a product's mutable fields.
type UpdateProductCommand struct {
	ProductID   string
	Name        string
	Description string
	Category    string
	Unit        string
	WeightKg    float64
}
