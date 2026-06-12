package model

import (
	"time"
)

// CustomerStatus represents the status of a customer.
type CustomerStatus string

const (
	CustomerStatusPending   CustomerStatus = "pending"
	CustomerStatusActive    CustomerStatus = "active"
	CustomerStatusRejected  CustomerStatus = "rejected"
	CustomerStatusSuspended CustomerStatus = "suspended"
)

// String returns the string representation.
func (s CustomerStatus) String() string {
	return string(s)
}

// IsValid checks if the status is valid.
func (s CustomerStatus) IsValid() bool {
	switch s {
	case CustomerStatusPending, CustomerStatusActive, CustomerStatusRejected, CustomerStatusSuspended:
		return true
	default:
		return false
	}
}

// Customer represents a customer (company) in the system.
type Customer struct {
	ID             string
	CompanyName    string
	Email          string
	Phone          string
	VatNumber      string
	Address        string
	City           string
	PostalCode     string
	Country        string
	Status         CustomerStatus
	CompanyAdminID string
	RejectReason   string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
