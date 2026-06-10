package validator_test

import (
	"testing"

	"github.com/beabys/wms/auth-service/internal/application/auth/validator"
	"github.com/stretchr/testify/assert"
)

func TestEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "test@example.com", false},
		{"valid email with plus", "test+tag@example.com", false},
		{"valid email with dots", "test.name@example.co.uk", false},
		{"empty email", "", true},
		{"invalid email no domain", "test@", true},
		{"invalid email no at", "testexample.com", true},
		{"invalid email spaces", "test @example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Email(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid password long enough", "password123", false},
		{"valid password exactly 8 chars", "12345678", false},
		{"empty password", "", true},
		{"too short 7 chars", "1234567", true},
		{"too short 1 char", "a", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Password(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{"non-empty token", "some-token", false},
		{"empty token", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Token(tt.token)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
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
		{"non-empty invite token", "invite-token", false},
		{"empty invite token", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.InviteToken(tt.token)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"valid UUID", "550e8400-e29b-41d4-a716-446655440000", false},
		{"empty id", "", true},
		{"invalid UUID format", "not-a-uuid", true},
		{"short string", "abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.UserID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
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
	}{
		{"non-empty field", "value", "name", false},
		{"empty field", "", "name", true},
		{"empty field with spaces", "   ", "name", false}, // NotEmpty doesn't trim
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.NotEmpty(tt.field, tt.fieldName)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.fieldName)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
