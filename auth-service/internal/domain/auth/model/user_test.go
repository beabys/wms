package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPasswordHash(t *testing.T) {
	t.Run("valid password creates hash", func(t *testing.T) {
		hash, err := NewPasswordHash("password123")
		assert.NoError(t, err)
		assert.NotEmpty(t, string(hash))
	})

	t.Run("different passwords produce different hashes", func(t *testing.T) {
		h1, _ := NewPasswordHash("password123")
		h2, _ := NewPasswordHash("password456")
		assert.NotEqual(t, string(h1), string(h2))
	})
}

func TestPasswordHashVerify(t *testing.T) {
	t.Run("correct password verifies", func(t *testing.T) {
		hash, err := NewPasswordHash("secret123")
		require.NoError(t, err)
		assert.True(t, hash.Verify("secret123"))
	})

	t.Run("wrong password does not verify", func(t *testing.T) {
		hash, err := NewPasswordHash("secret123")
		require.NoError(t, err)
		assert.False(t, hash.Verify("wrongpass"))
	})

	t.Run("empty password does not verify", func(t *testing.T) {
		hash, err := NewPasswordHash("secret123")
		require.NoError(t, err)
		assert.False(t, hash.Verify(""))
	})
}

func TestNewUser(t *testing.T) {
	t.Run("valid user creates successfully", func(t *testing.T) {
		user, err := NewUser("test@example.com", "password123", "Test User", "admin", nil)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.NotEmpty(t, user.ID)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "Test User", user.Name)
		assert.Equal(t, "admin", user.Role)
		assert.Nil(t, user.CustomerID)
		assert.True(t, user.Active)
		assert.False(t, user.CreatedAt.IsZero())
		assert.False(t, user.UpdatedAt.IsZero())
	})

	t.Run("user with customer ID creates successfully", func(t *testing.T) {
		cid := "cust-123"
		user, err := NewUser("user@customer.com", "password123", "Customer User", "customer", &cid)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.NotNil(t, user.CustomerID)
		assert.Equal(t, "cust-123", *user.CustomerID)
	})

	t.Run("invalid email format returns error", func(t *testing.T) {
		_, err := NewUser("not-an-email", "password123", "Test", "admin", nil)
		assert.Error(t, err)
	})
}

func TestUserUpdatePassword(t *testing.T) {
	t.Run("valid password update succeeds", func(t *testing.T) {
		user, err := NewUser("test@example.com", "password123", "Test", "admin", nil)
		require.NoError(t, err)

		oldHash := user.Password
		err = user.UpdatePassword("newpassword456")
		assert.NoError(t, err)
		assert.NotEqual(t, string(oldHash), string(user.Password))
		assert.True(t, user.Password.Verify("newpassword456"))
	})
}

func TestUserDeactivate(t *testing.T) {
	t.Run("deactivate sets active to false", func(t *testing.T) {
		user, _ := NewUser("test@example.com", "password123", "Test", "admin", nil)
		assert.True(t, user.Active)

		user.Deactivate()
		assert.False(t, user.Active)
	})
}
