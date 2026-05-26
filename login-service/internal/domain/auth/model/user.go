package model

import "time"

// UserStatus represents the status of a user account.
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// User is the root entity for authentication.
type User struct {
	ID           string
	Email        Email
	PasswordHash PasswordHash
	Role         Role
	CustomerID   string
	Status       UserStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// IsActive checks if the user account is active.
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// CanLogin checks if the user can log in based on status.
func (u *User) CanLogin() bool {
	return u.Status == UserStatusActive
}
