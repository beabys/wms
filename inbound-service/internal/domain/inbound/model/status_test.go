package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name    string
		current Status
		target  Status
		want    bool
	}{
		// submitted transitions
		{"submitted -> inspected", StatusSubmitted, StatusInspected, true},
		{"submitted -> approved", StatusSubmitted, StatusApproved, false},
		{"submitted -> flagged", StatusSubmitted, StatusFlagged, false},
		{"submitted -> held", StatusSubmitted, StatusHeld, false},
		{"submitted -> released", StatusSubmitted, StatusReleased, false},
		// inspected transitions
		{"inspected -> approved", StatusInspected, StatusApproved, true},
		{"inspected -> flagged", StatusInspected, StatusFlagged, true},
		{"inspected -> held", StatusInspected, StatusHeld, true},
		{"inspected -> submitted", StatusInspected, StatusSubmitted, false},
		{"inspected -> released", StatusInspected, StatusReleased, false},
		// approved is terminal
		{"approved -> released", StatusApproved, StatusReleased, false},
		{"approved -> flagged", StatusApproved, StatusFlagged, false},
		{"approved -> held", StatusApproved, StatusHeld, false},
		{"approved -> submitted", StatusApproved, StatusSubmitted, false},
		// flagged is terminal
		{"flagged -> released", StatusFlagged, StatusReleased, false},
		{"flagged -> approved", StatusFlagged, StatusApproved, false},
		// held -> released
		{"held -> released", StatusHeld, StatusReleased, true},
		{"held -> approved", StatusHeld, StatusApproved, false},
		{"held -> flagged", StatusHeld, StatusFlagged, false},
		// released is terminal
		{"released -> held", StatusReleased, StatusHeld, false},
		{"released -> submitted", StatusReleased, StatusSubmitted, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.current.CanTransitionTo(tt.target)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateTransition(t *testing.T) {
	t.Run("valid transition", func(t *testing.T) {
		err := ValidateTransition(StatusSubmitted, StatusInspected)
		assert.NoError(t, err)
	})

	t.Run("invalid transition returns error", func(t *testing.T) {
		err := ValidateTransition(StatusSubmitted, StatusApproved)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status transition")
	})
}
