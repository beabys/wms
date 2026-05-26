package command

import "time"

// InspectInboundCommand records an inspection.
type InspectInboundCommand struct {
	InboundID   string
	InspectorID string
	Notes       string
	Photos      []string
	InspectedAt time.Time
	Passed      bool
}
