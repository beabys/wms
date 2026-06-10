package model

import "time"

// CompanyRole represents a role with permissions within a company.
type CompanyRole struct {
	ID          string
	CustomerID  string
	Name        string
	Permissions []string
	IsDefault   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
