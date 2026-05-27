package model

import (
	"fmt"
	"strings"
	"time"
)

// Product is the aggregate root for a product in inventory.
type Product struct {
	ID                string
	SKU               string
	Name              string
	Description       string
	Category          string
	Unit              string
	WeightKg          float64
	IsActive          bool
	LowStockThreshold float64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// RegisterProduct creates a new active product with the given attributes.
func RegisterProduct(id, sku, name, description, category, unit string, weightKg, lowStockThreshold float64) (*Product, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("product id is required")
	}
	if strings.TrimSpace(sku) == "" {
		return nil, fmt.Errorf("sku is required")
	}
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name is required")
	}
	if weightKg < 0 {
		return nil, fmt.Errorf("weight cannot be negative")
	}
	if lowStockThreshold < 0 {
		return nil, fmt.Errorf("low stock threshold cannot be negative")
	}

	now := time.Now()
	return &Product{
		ID:                id,
		SKU:               sku,
		Name:              name,
		Description:       description,
		Category:          category,
		Unit:              unit,
		WeightKg:          weightKg,
		IsActive:          true,
		LowStockThreshold: lowStockThreshold,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

// Update modifies mutable product attributes.
func (p *Product) Update(name, description, category, unit string, weightKg float64) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name is required")
	}
	if weightKg < 0 {
		return fmt.Errorf("weight cannot be negative")
	}

	p.Name = name
	p.Description = description
	p.Category = category
	p.Unit = unit
	p.WeightKg = weightKg
	p.UpdatedAt = time.Now()
	return nil
}

// Archive sets the product as inactive.
func (p *Product) Archive() {
	p.IsActive = false
	p.UpdatedAt = time.Now()
}
