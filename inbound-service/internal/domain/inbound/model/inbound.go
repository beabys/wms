package model

import (
	"fmt"
	"strings"
	"time"
)

// Inbound is the aggregate root for an inbound shipment.
type Inbound struct {
	ID                 string
	CustomerID         string
	Status             Status
	ExpectedDate       string
	PackagingProfileID string
	Notes              string
	CreatedAt          time.Time
	UpdatedAt          time.Time

	Items      []InboundItem
	Inspection *Inspection
	HoldRecord *Hold
}

// SubmitInbound creates a new Inbound with submitted status.
func SubmitInbound(id, customerID, expectedDate, notes string, items []InboundItem) (*Inbound, error) {
	if id == "" {
		return nil, fmt.Errorf("inbound id is required")
	}
	if customerID == "" {
		return nil, fmt.Errorf("customer id is required")
	}
	if expectedDate == "" {
		return nil, fmt.Errorf("expected date is required")
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("at least one item is required")
	}
	for i, item := range items {
		if err := item.Validate(); err != nil {
			return nil, fmt.Errorf("item[%d]: %w", i, err)
		}
	}

	return &Inbound{
		ID:           id,
		CustomerID:   customerID,
		Status:       StatusSubmitted,
		ExpectedDate: expectedDate,
		Notes:        notes,
		Items:        items,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// Inspect records an inspection and transitions status to inspected.
func (in *Inbound) Inspect(inspection Inspection) error {
	if err := ValidateTransition(in.Status, StatusInspected); err != nil {
		return err
	}
	if err := inspection.Validate(); err != nil {
		return err
	}
	in.Inspection = &inspection
	in.Status = StatusInspected
	in.UpdatedAt = time.Now()
	return nil
}

// Approve transitions status to approved.
func (in *Inbound) Approve() error {
	if err := ValidateTransition(in.Status, StatusApproved); err != nil {
		return err
	}
	in.Status = StatusApproved
	in.UpdatedAt = time.Now()
	return nil
}

// Flag transitions status to flagged.
func (in *Inbound) Flag(reason string) error {
	if err := ValidateTransition(in.Status, StatusFlagged); err != nil {
		return err
	}
	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("flag reason is required")
	}
	in.Status = StatusFlagged
	in.UpdatedAt = time.Now()
	return nil
}

// PlaceHold transitions status to held.
func (in *Inbound) PlaceHold(reason string) error {
	if err := ValidateTransition(in.Status, StatusHeld); err != nil {
		return err
	}
	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("hold reason is required")
	}

	now := time.Now()
	in.HoldRecord = &Hold{
		Reason:    reason,
		CreatedAt: now,
	}
	in.Status = StatusHeld
	in.UpdatedAt = now
	return nil
}

// Release transitions status from held to released.
func (in *Inbound) Release(releasedBy string) error {
	if err := ValidateTransition(in.Status, StatusReleased); err != nil {
		return err
	}
	if strings.TrimSpace(releasedBy) == "" {
		return fmt.Errorf("released by is required")
	}
	if in.HoldRecord == nil {
		return fmt.Errorf("no hold record to release")
	}

	now := time.Now()
	in.HoldRecord.ReleasedAt = &now
	in.HoldRecord.ReleasedBy = releasedBy
	in.Status = StatusReleased
	in.UpdatedAt = now
	return nil
}
