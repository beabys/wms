package validator

import (
	"fmt"
	"net/mail"

	"github.com/google/uuid"
)

// Email validates email format.
func Email(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email format: %w", err)
	}
	return nil
}

// Password validates password meets minimum length.
func Password(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}

// Token validates a token string is non-empty.
func Token(token string) error {
	if token == "" {
		return fmt.Errorf("token is required")
	}
	return nil
}

// InviteToken validates an invite token string is non-empty.
func InviteToken(token string) error {
	if token == "" {
		return fmt.Errorf("invite token is required")
	}
	return nil
}

// UserID validates user ID format (UUID).
func UserID(id string) error {
	if id == "" {
		return fmt.Errorf("user ID is required")
	}
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}
	return nil
}

// NotEmpty validates a field is non-empty.
func NotEmpty(field, name string) error {
	if field == "" {
		return fmt.Errorf("%s is required", name)
	}
	return nil
}


