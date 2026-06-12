package model

import "time"

// AuditAction represents the action performed on a customer.
type AuditAction string

const (
	AuditActionApproved AuditAction = "approved"
	AuditActionRejected AuditAction = "rejected"
	AuditActionSuspended AuditAction = "suspended"
	AuditActionRestored  AuditAction = "restored"
	AuditActionUpdated   AuditAction = "updated"
)

// String returns the string representation.
func (a AuditAction) String() string {
	return string(a)
}

// AuditLog represents a single audit entry for a customer.
type AuditLog struct {
	ID          string
	CustomerID  string
	Action      AuditAction
	PerformedBy string
	Details     string
	CreatedAt   time.Time
}
