package validator

import (
	"fmt"
	"net/mail"
)

// ValidateEmail checks email format.
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidatePassword checks password requirements.
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}

// ValidateCreateUser checks all required fields for user creation.
func ValidateCreateUser(email, password, name, role string) error {
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if err := ValidatePassword(password); err != nil {
		return err
	}
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if role == "" {
		return fmt.Errorf("role is required")
	}
	return nil
}

// ValidateInviteEmail checks the email for invite.
func ValidateInviteEmail(email string) error {
	return ValidateEmail(email)
}
