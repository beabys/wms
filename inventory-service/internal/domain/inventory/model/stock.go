package model

import (
	"fmt"
	"strings"
	"time"
)

// StockEntry represents a quantity of a product at a bin location.
type StockEntry struct {
	ID              string
	ProductID       string
	BinLocationID   string
	Quantity        float64
	ReservedQuantity float64
	LotNumber       string
	ExpiryDate      *time.Time
	Status          string
	LastUpdated     time.Time
}

// NewStockEntry creates a new stock entry.
func NewStockEntry(id, productID, binLocationID string, quantity float64, lotNumber string, expiryDate *time.Time) (*StockEntry, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("stock entry id is required")
	}
	if strings.TrimSpace(productID) == "" {
		return nil, fmt.Errorf("product id is required")
	}
	if quantity < 0 {
		return nil, fmt.Errorf("quantity cannot be negative")
	}

	return &StockEntry{
		ID:               id,
		ProductID:        productID,
		BinLocationID:    binLocationID,
		Quantity:         quantity,
		ReservedQuantity: 0,
		LotNumber:        lotNumber,
		ExpiryDate:       expiryDate,
		Status:           "active",
		LastUpdated:      time.Now(),
	}, nil
}

// Adjust changes the quantity by delta (positive or negative).
func (s *StockEntry) Adjust(delta float64, reason string) error {
	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("adjust reason is required")
	}
	if s.Quantity+delta < 0 {
		return fmt.Errorf("insufficient stock: have %f, need %f", s.Quantity, -delta)
	}
	s.Quantity += delta
	s.LastUpdated = time.Now()
	return nil
}

// Reserve marks a quantity as reserved.
func (s *StockEntry) Reserve(qty float64) error {
	if qty <= 0 {
		return fmt.Errorf("reserve quantity must be positive")
	}
	available := s.Quantity - s.ReservedQuantity
	if qty > available {
		return fmt.Errorf("insufficient available stock: have %f, need %f", available, qty)
	}
	s.ReservedQuantity += qty
	s.LastUpdated = time.Now()
	return nil
}

// Release reduces the reserved quantity.
func (s *StockEntry) Release(qty float64) error {
	if qty <= 0 {
		return fmt.Errorf("release quantity must be positive")
	}
	if qty > s.ReservedQuantity {
		return fmt.Errorf("cannot release %f, only %f reserved", qty, s.ReservedQuantity)
	}
	s.ReservedQuantity -= qty
	s.LastUpdated = time.Now()
	return nil
}

// IsLowStock checks if available stock is below the threshold.
func (s *StockEntry) IsLowStock(threshold float64) bool {
	return (s.Quantity - s.ReservedQuantity) < threshold
}
