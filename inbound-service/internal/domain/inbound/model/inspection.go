package model

import (
	"fmt"
	"time"
)

// Inspection records the result of inspecting an inbound shipment.
type Inspection struct {
	InspectorID string
	Notes       string
	Photos      []string
	InspectedAt time.Time
	Passed      bool
}

// Validate checks inspection fields.
func (i Inspection) Validate() error {
	if i.InspectorID == "" {
		return fmt.Errorf("inspector_id is required")
	}
	if i.InspectedAt.IsZero() {
		return fmt.Errorf("inspected_at is required")
	}
	return nil
}
