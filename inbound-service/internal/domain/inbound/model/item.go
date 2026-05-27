package model

import "fmt"

// InboundItem represents a single SKU line within an inbound shipment.
type InboundItem struct {
	ID               string
	InboundID        string
	SKU              string
	QuantityDeclared int32
	QuantityReceived int32
	Dimensions       string // JSON string: {"x":10,"y":20,"z":30}
	Weight           float64
}

// Validate checks item fields.
func (i InboundItem) Validate() error {
	if i.SKU == "" {
		return fmt.Errorf("sku is required")
	}
	if i.QuantityDeclared <= 0 {
		return fmt.Errorf("quantity_declared must be positive")
	}
	if i.QuantityReceived < 0 {
		return fmt.Errorf("quantity_received cannot be negative")
	}
	return nil
}
