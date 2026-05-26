package model

import "fmt"

// Status represents the lifecycle state of an inbound.
type Status string

const (
	StatusSubmitted  Status = "submitted"
	StatusInspected  Status = "inspected"
	StatusApproved   Status = "approved"
	StatusFlagged    Status = "flagged"
	StatusHeld       Status = "held"
	StatusReleased   Status = "released"
)

// ValidTransitions defines allowed status transitions.
var ValidTransitions = map[Status][]Status{
	StatusSubmitted: {StatusInspected},
	StatusInspected: {StatusApproved, StatusFlagged, StatusHeld},
	StatusApproved:  {},
	StatusFlagged:   {},
	StatusHeld:      {StatusReleased},
	StatusReleased:  {},
}

// CanTransitionTo returns true if the current status allows moving to target.
func (s Status) CanTransitionTo(target Status) bool {
	allowed, ok := ValidTransitions[s]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == target {
			return true
		}
	}
	return false
}

// ValidateTransition returns an error if transitioning from current to target is invalid.
func ValidateTransition(current, target Status) error {
	if !current.CanTransitionTo(target) {
		return fmt.Errorf("invalid status transition: %s -> %s", current, target)
	}
	return nil
}
