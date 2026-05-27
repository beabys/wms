package model

import (
	"fmt"
	"net/mail"
)

// Email is a value object representing an email address.
type Email string

// NewEmail creates a new Email value object with validation.
func NewEmail(raw string) (Email, error) {
	if raw == "" {
		return "", fmt.Errorf("email cannot be empty")
	}
	addr, err := mail.ParseAddress(raw)
	if err != nil {
		return "", fmt.Errorf("invalid email format: %w", err)
	}
	return Email(addr.Address), nil
}

// String returns the email string.
func (e Email) String() string {
	return string(e)
}
