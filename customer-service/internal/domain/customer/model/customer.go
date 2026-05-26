package model

import (
	"fmt"
	"time"
)

// CustomerStatus represents the lifecycle state of a customer.
type CustomerStatus string

const (
	CustomerStatusPending   CustomerStatus = "pending"
	CustomerStatusActive    CustomerStatus = "active"
	CustomerStatusSuspended CustomerStatus = "suspended"
)

// Customer is the aggregate root for the customer bounded context.
type Customer struct {
	ID            string
	CompanyName   string
	VATNumber     string
	Address       Address
	Status        CustomerStatus
	RateCardID    string
	CreditBalance int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Approve transitions customer from pending to active.
func (c *Customer) Approve() error {
	if c.Status != CustomerStatusPending {
		return fmt.Errorf("cannot approve customer in status %s", c.Status)
	}
	c.Status = CustomerStatusActive
	c.UpdatedAt = time.Now()
	return nil
}

// Suspend transitions customer from active to suspended.
func (c *Customer) Suspend() error {
	if c.Status != CustomerStatusActive {
		return fmt.Errorf("cannot suspend customer in status %s", c.Status)
	}
	c.Status = CustomerStatusSuspended
	c.UpdatedAt = time.Now()
	return nil
}

// Reactivate transitions customer from suspended back to active.
func (c *Customer) Reactivate() error {
	if c.Status != CustomerStatusSuspended {
		return fmt.Errorf("cannot reactivate customer in status %s", c.Status)
	}
	c.Status = CustomerStatusActive
	c.UpdatedAt = time.Now()
	return nil
}
