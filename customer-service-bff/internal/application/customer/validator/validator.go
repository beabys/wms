package validator

import (
	"fmt"
	"net/mail"
)

// ValidateRegisterCustomer checks all required fields for customer registration.
func ValidateRegisterCustomer(token, companyName, email, password string) error {
	if token == "" {
		return fmt.Errorf("token is required")
	}
	if companyName == "" {
		return fmt.Errorf("company_name is required")
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

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

// ValidateAssignCompanyRole checks required fields for assigning a company role.
func ValidateAssignCompanyRole(customerID, userID, roleName string) error {
	if customerID == "" {
		return fmt.Errorf("customer_id is required")
	}
	if userID == "" {
		return fmt.Errorf("user_id is required")
	}
	if roleName == "" {
		return fmt.Errorf("role_name is required")
	}
	return nil
}

// ValidateCustomerID checks that a customer ID is a valid UUID (basic length check).
func ValidateCustomerID(id string) error {
	if id == "" {
		return fmt.Errorf("customer id is required")
	}
	if len(id) != 36 {
		return fmt.Errorf("invalid customer id format")
	}
	return nil
}
