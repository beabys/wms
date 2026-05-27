package model

import "time"

// Hold represents a hold placed on an inbound.
type Hold struct {
	ID         string
	Reason     string
	CreatedAt  time.Time
	ReleasedAt *time.Time // nil until released
	ReleasedBy string
}
