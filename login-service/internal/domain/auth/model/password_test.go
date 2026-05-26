package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPasswordHash(t *testing.T) {
	tests := []struct {
		name    string
		pass    string
		wantErr bool
	}{
		{"valid password", "securePass123!", false},
		{"minimum length", "12345678", false},
		{"empty", "", true},
		{"too short", "1234567", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := NewPasswordHash(tt.pass)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, hash.String())
			assert.NotEqual(t, tt.pass, hash.String())
		})
	}
}

func TestPasswordHash_Verify(t *testing.T) {
	hash, err := NewPasswordHash("mySecretPassword123")
	require.NoError(t, err)

	assert.True(t, hash.Verify("mySecretPassword123"))
	assert.False(t, hash.Verify("wrongPassword"))
	assert.False(t, hash.Verify(""))
}

func TestPasswordHashFromHash(t *testing.T) {
	original := "mySecretPassword!"
	hash, err := NewPasswordHash(original)
	require.NoError(t, err)

	// Create from existing hash
	hash2 := PasswordHashFromHash(hash.String())
	assert.True(t, hash2.Verify(original))
	assert.Equal(t, hash.String(), hash2.String())
}
