package validator

import (
	"fmt"
	"net/mail"
)

// Email validates an email address format.
func Email(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format: %s", email)
	}
	return nil
}

// NotEmpty checks that a field is not empty.
func NotEmpty(field, name string) error {
	if field == "" {
		return fmt.Errorf("%s is required", name)
	}
	return nil
}

// CustomerID validates a customer ID is not empty.
func CustomerID(id string) error {
	if id == "" {
		return fmt.Errorf("customer ID is required")
	}
	return nil
}

// CompanyName validates a company name is not empty.
func CompanyName(name string) error {
	return NotEmpty(name, "company name")
}

// InviteToken validates an invite token is not empty.
func InviteToken(token string) error {
	return NotEmpty(token, "invite token")
}
