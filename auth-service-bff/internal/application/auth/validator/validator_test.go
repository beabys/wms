package validator

import (
	"strings"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "user@example.com", false},
		{"valid with plus", "user+tag@example.com", false},
		{"valid subdomain", "user@sub.example.com", false},
		{"empty email", "", true},
		{"no @", "userexample.com", true},
		{"no domain", "user@", true},
		{"no local", "@example.com", true},
		{"double dot", "user@example..com", true},
		{"spaces", "user @example.com", true},
		{"missing domain", "user@", true},
		{"no at sign", "userexample.com", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateEmail_ErrorMessages(t *testing.T) {
	if err := ValidateEmail(""); err == nil {
		t.Fatal("expected error for empty email")
	} else if err.Error() != "email is required" {
		t.Errorf("expected 'email is required', got %q", err.Error())
	}

	if err := ValidateEmail("invalid"); err == nil {
		t.Fatal("expected error for invalid email")
	} else if err.Error() != "invalid email format" {
		t.Errorf("expected 'invalid email format', got %q", err.Error())
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid password", "password123", false},
		{"valid long", strings.Repeat("a", 100), false},
		{"empty", "", true},
		{"too short", "short1", true},
		{"exactly 7 chars", "1234567", true},
		{"exactly 8 chars", "12345678", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidatePassword_ErrorMessages(t *testing.T) {
	if err := ValidatePassword(""); err == nil {
		t.Fatal("expected error for empty password")
	} else if err.Error() != "password is required" {
		t.Errorf("expected 'password is required', got %q", err.Error())
	}

	if err := ValidatePassword("abc"); err == nil {
		t.Fatal("expected error for short password")
	} else if err.Error() != "password must be at least 8 characters" {
		t.Errorf("expected 'password must be at least 8 characters', got %q", err.Error())
	}
}

func TestValidateCreateUser(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		nameVal  string
		role     string
		wantErr  bool
	}{
		{"all valid", "user@example.com", "password123", "John", "admin", false},
		{"empty email", "", "password123", "John", "admin", true},
		{"invalid email", "not-an-email", "password123", "John", "admin", true},
		{"empty password", "user@example.com", "", "John", "admin", true},
		{"short password", "user@example.com", "short", "John", "admin", true},
		{"empty name", "user@example.com", "password123", "", "admin", true},
		{"empty role", "user@example.com", "password123", "John", "", true},
		{"multiple failures returns first", "", "", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateUser(tt.email, tt.password, tt.nameVal, tt.role)
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateInviteEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "invite@example.com", false},
		{"empty", "", true},
		{"invalid", "bad-email", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInviteEmail(tt.email)
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
