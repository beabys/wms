package model

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// PasswordHash is a bcrypt hash of a password.
type PasswordHash string

// NewPasswordHash creates a bcrypt hash from a plain-text password.
func NewPasswordHash(password string) (PasswordHash, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return PasswordHash(string(bytes)), nil
}

// Verify checks a plain-text password against the hash.
func (h PasswordHash) Verify(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(string(h)), []byte(password)) == nil
}

// User represents a user in the system.
type User struct {
	ID         string
	Email      string
	Password   PasswordHash
	Name       string
	Role       string       // e.g., "admin", "warehouse_staff", "billing_manager", "customer"
	CustomerID *string      // nil for internal system users
	Active     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewUser creates a new user. Input must be pre-validated by the validator package.
func NewUser(email, password, name, role string, customerID *string) (*User, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, fmt.Errorf("invalid email format: %w", err)
	}

	hash, err := NewPasswordHash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now()
	return &User{
		ID:         uuid.New().String(),
		Email:      email,
		Password:   hash,
		Name:       name,
		Role:       role,
		CustomerID: customerID,
		Active:     true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// UpdatePassword hashes and updates the user's password. Input must be pre-validated.
func (u *User) UpdatePassword(newPassword string) error {
	hash, err := NewPasswordHash(newPassword)
	if err != nil {
		return err
	}
	u.Password = hash
	u.UpdatedAt = time.Now()
	return nil
}

// Deactivate marks the user as inactive.
func (u *User) Deactivate() {
	u.Active = false
	u.UpdatedAt = time.Now()
}
