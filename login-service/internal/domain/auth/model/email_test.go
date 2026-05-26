package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid email", "user@example.com", false},
		{"valid email with plus", "user+tag@example.com", false},
		{"valid email subdomain", "user@sub.example.com", false},
		{"empty", "", true},
		{"no domain", "user", true},
		{"no at sign", "userexample.com", true},
		{"double at", "user@ex@ample.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.input, email.String())
		})
	}
}
