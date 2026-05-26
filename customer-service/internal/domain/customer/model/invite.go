package model

import (
	"fmt"
	"time"
)

// InviteLink represents an invitation sent to a prospective customer.
type InviteLink struct {
	ID        string
	Email     string
	Token     string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}

// IsExpired checks whether the invite has expired.
func (i *InviteLink) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

// MarkUsed marks the invite as used.
func (i *InviteLink) MarkUsed() error {
	if i.Used {
		return fmt.Errorf("invite %s already used", i.ID)
	}
	if i.IsExpired() {
		return fmt.Errorf("invite %s has expired", i.ID)
	}
	i.Used = true
	return nil
}
