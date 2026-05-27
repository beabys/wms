package model

import (
	"fmt"
	"strings"
	"time"
)

// BinLocation represents a physical storage location in the warehouse.
type BinLocation struct {
	ID            string
	WarehouseZone string
	Aisle         string
	Rack          string
	Shelf         string
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewBinLocation creates a new active bin location.
func NewBinLocation(id, zone, aisle, rack, shelf string) (*BinLocation, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("bin location id is required")
	}
	if strings.TrimSpace(zone) == "" {
		return nil, fmt.Errorf("warehouse zone is required")
	}
	if strings.TrimSpace(aisle) == "" {
		return nil, fmt.Errorf("aisle is required")
	}
	if strings.TrimSpace(rack) == "" {
		return nil, fmt.Errorf("rack is required")
	}
	if strings.TrimSpace(shelf) == "" {
		return nil, fmt.Errorf("shelf is required")
	}

	now := time.Now()
	return &BinLocation{
		ID:            id,
		WarehouseZone: zone,
		Aisle:         aisle,
		Rack:          rack,
		Shelf:         shelf,
		IsActive:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}
