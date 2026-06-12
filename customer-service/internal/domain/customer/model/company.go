package model

import "time"

// Company represents a company entity owned by a customer.
type Company struct {
	ID         string
	CustomerID string
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
