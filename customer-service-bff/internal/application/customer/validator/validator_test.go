package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRegisterCustomer(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		companyName string
		email       string
		password    string
		wantErr     bool
		errMsg      string
	}{
		{"all valid", "token-123", "ACME Corp", "admin@acme.com", "secure-pass", false, ""},
		{"missing token", "", "ACME Corp", "admin@acme.com", "secure-pass", true, "token is required"},
		{"missing company_name", "token-123", "", "admin@acme.com", "secure-pass", true, "company_name is required"},
		{"missing email", "token-123", "ACME Corp", "", "secure-pass", true, "email is required"},
		{"invalid email", "token-123", "ACME Corp", "not-an-email", "secure-pass", true, "invalid email format"},
		{"missing password", "token-123", "ACME Corp", "admin@acme.com", "", true, "password is required"},
		{"all missing", "", "", "", "", true, "token is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRegisterCustomer(tt.token, tt.companyName, tt.email, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "user@example.com", false},
		{"empty", "", true},
		{"invalid", "not-an-email", true},
		{"missing domain", "user@", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateAssignCompanyRole(t *testing.T) {
	tests := []struct {
		name       string
		customerID string
		userID     string
		roleName   string
		wantErr    bool
		errMsg     string
	}{
		{"all valid", "cust-1", "user-1", "admin", false, ""},
		{"missing customer_id", "", "user-1", "admin", true, "customer_id is required"},
		{"missing user_id", "cust-1", "", "admin", true, "user_id is required"},
		{"missing role_name", "cust-1", "user-1", "", true, "role_name is required"},
		{"all missing", "", "", "", true, "customer_id is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAssignCompanyRole(tt.customerID, tt.userID, tt.roleName)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCustomerID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		wantErr bool
	}{
		{"valid UUID format", "550e8400-e29b-41d4-a716-446655440000", false},
		{"empty", "", true},
		{"too short", "abc", true},
		{"wrong format", "123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCustomerID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
