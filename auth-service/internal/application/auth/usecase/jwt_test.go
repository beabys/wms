package usecase

import (
	"testing"
	"time"

	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestService(t *testing.T) *Service {
	t.Helper()
	privPEM, pubPEM, err := GenerateKeyPair()
	require.NoError(t, err)

	svc, err := NewService(privPEM, pubPEM, 15*time.Minute, 7*24*time.Hour)
	require.NoError(t, err)
	return svc
}

func TestNewService(t *testing.T) {
	t.Run("valid keys create service", func(t *testing.T) {
		svc := setupTestService(t)
		assert.NotNil(t, svc)
	})

	t.Run("invalid private key returns error", func(t *testing.T) {
		_, err := NewService("invalid-pem", "invalid-pem", 15*time.Minute, 7*24*time.Hour)
		assert.Error(t, err)
	})
}

func TestGenerateAccessToken(t *testing.T) {
	t.Run("generates valid token", func(t *testing.T) {
		svc := setupTestService(t)
		user := &model.User{ID: "user-1", Role: "admin", Email: "test@example.com", Name: "Test"}

		token, err := svc.GenerateAccessToken(user)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("token with customerID includes it in claims", func(t *testing.T) {
		svc := setupTestService(t)
		cid := "cust-123"
		user := &model.User{ID: "user-2", Role: "customer", Email: "cust@test.com", Name: "Cust", CustomerID: &cid}

		token, err := svc.GenerateAccessToken(user)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestValidateToken(t *testing.T) {
	t.Run("valid token returns claims", func(t *testing.T) {
		svc := setupTestService(t)
		user := &model.User{ID: "user-1", Role: "admin", Email: "test@example.com", Name: "Test"}

		token, err := svc.GenerateAccessToken(user)
		require.NoError(t, err)

		claims, err := svc.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, "user-1", claims.UserID)
		assert.Equal(t, "admin", claims.Role)
	})

	t.Run("expired token returns error", func(t *testing.T) {
		privPEM, pubPEM, err := GenerateKeyPair()
		require.NoError(t, err)

		// Create service with -1 minute TTL so token is immediately expired
		svc, err := NewService(privPEM, pubPEM, -1*time.Minute, 7*24*time.Hour)
		require.NoError(t, err)

		user := &model.User{ID: "user-1", Role: "admin", Email: "test@example.com", Name: "Test"}

		token, err := svc.GenerateAccessToken(user)
		require.NoError(t, err)

		_, err = svc.ValidateToken(token)
		assert.Error(t, err)
	})

	t.Run("invalid token string returns error", func(t *testing.T) {
		svc := setupTestService(t)
		_, err := svc.ValidateToken("invalid.token.string")
		assert.Error(t, err)
	})

	t.Run("token from different key is rejected", func(t *testing.T) {
		svc1 := setupTestService(t)
		svc2 := setupTestService(t)

		user := &model.User{ID: "user-1", Role: "admin"}
		token, err := svc1.GenerateAccessToken(user)
		require.NoError(t, err)

		_, err = svc2.ValidateToken(token)
		assert.Error(t, err)
	})
}

func TestGenerateRefreshToken(t *testing.T) {
	t.Run("generates unique tokens", func(t *testing.T) {
		svc := setupTestService(t)

		t1, err := svc.GenerateRefreshToken()
		require.NoError(t, err)
		assert.NotEmpty(t, t1)

		t2, err := svc.GenerateRefreshToken()
		require.NoError(t, err)
		assert.NotEqual(t, t1, t2)
	})
}

func TestGetPublicKeyPEM(t *testing.T) {
	t.Run("returns valid PEM", func(t *testing.T) {
		svc := setupTestService(t)

		pem, err := svc.GetPublicKeyPEM()
		require.NoError(t, err)
		assert.Contains(t, pem, "-----BEGIN PUBLIC KEY-----")
		assert.Contains(t, pem, "-----END PUBLIC KEY-----")
	})
}

func TestGetAccessTokenTTL(t *testing.T) {
	t.Run("returns configured TTL", func(t *testing.T) {
		svc := setupTestService(t)
		ttl := svc.GetAccessTokenTTL()
		assert.Equal(t, int64(900), ttl) // 15 minutes
	})
}

func TestGenerateKeyPair(t *testing.T) {
	t.Run("generates valid PEM keys", func(t *testing.T) {
		priv, pub, err := GenerateKeyPair()
		require.NoError(t, err)
		assert.Contains(t, priv, "-----BEGIN PRIVATE KEY-----")
		assert.Contains(t, pub, "-----BEGIN PUBLIC KEY-----")

		// Verify they can be used to create a service
		_, err = NewService(priv, pub, 15*time.Minute, 7*24*time.Hour)
		require.NoError(t, err)
	})
}
