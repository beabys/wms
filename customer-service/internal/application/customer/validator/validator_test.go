package validator

import (
	"testing"
)

func TestEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "user@example.com", false},
		{"valid email with subdomain", "user@sub.example.com", false},
		{"valid email with plus", "user+tag@example.com", false},
		{"empty email", "", true},
		{"invalid - no domain", "user@", true},
		{"invalid - no at sign", "userexample.com", true},
		{"invalid - spaces", "user @example.com", true},
		{"invalid - double dots", "user@example..com", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Email(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("Email() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotEmpty(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		fieldName string
		wantErr bool
		errMsg  string
	}{
		{"non-empty field", "value", "test field", false, ""},
		{"empty field", "", "test field", true, "test field is required"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NotEmpty(tt.field, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NotEmpty() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("NotEmpty() error message = %q, want %q", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestCompanyName(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"non-empty name", "Test Corp", false},
		{"empty name", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CompanyName(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompanyName() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestInviteToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{"non-empty token", "valid-token-123", false},
		{"empty token", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InviteToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("InviteToken() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestCustomerID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"non-empty id", "cust-123", false},
		{"empty id", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CustomerID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomerID() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
