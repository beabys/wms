//go:build integration

package repository

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/beabys/wms/pkg/database"
	"github.com/beabys/wms/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dbmocks "github.com/beabys/wms/auth-service/mocks/database"
)

var authTestUserID = uuid.New().String()

func setupAuthRepoTest(t *testing.T) (*AuthRepository, func()) {
	t.Helper()
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPass := getEnvOrDefault("DB_PASSWORD", "postgres")
	dbName := getEnvOrDefault("DB_NAME", "login_service_test")
	dbSSL := getEnvOrDefault("DB_SSLMODE", "disable")

	connStr := "postgres://" + dbUser + ":" + dbPass + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=" + dbSSL

	db, err := sqlx.Connect("pgx", connStr)
	require.NoError(t, err)

	// Clean up
	_, err = db.Exec("DELETE FROM refresh_tokens")
	require.NoError(t, err)
	_, err = db.Exec("DELETE FROM users")
	require.NoError(t, err)

	// Create a test user
	_, err = db.Exec(`INSERT INTO users (id, email, password_hash, name, role, active, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, true, NOW(), NOW())`,
		authTestUserID, "authtest@test.com", "hash", "Auth Test", "admin",
	)
	require.NoError(t, err)

	log := logger.NewSlogLogger(slog.LevelDebug)
	pg := &database.Postgres{DB: db}
	d := dbmocks.NewDatabase(t)
	d.On("Connect").Return(nil)
	d.On("Ping").Return(nil)
	d.On("Close").Return(nil)
	d.On("GetDBImpl").Return(pg)

	repo := NewAuthRepository(log, d)
	cleanup := func() {
		db.Exec("DELETE FROM refresh_tokens")
		db.Exec("DELETE FROM users")
		db.Close()
	}

	return repo, cleanup
}

func TestAuthRepository_StoreRefreshToken(t *testing.T) {
	repo, cleanup := setupAuthRepoTest(t)
	defer cleanup()

	t.Run("store token succeeds", func(t *testing.T) {
		hash := sha256Hex("test-refresh-token")
		err := repo.StoreRefreshToken(context.Background(), authTestUserID, hash, time.Now().Add(1*time.Hour))
		require.NoError(t, err)
	})
}

func TestAuthRepository_ValidateRefreshToken(t *testing.T) {
	repo, cleanup := setupAuthRepoTest(t)
	defer cleanup()

	t.Run("valid token returns user ID", func(t *testing.T) {
		hash := sha256Hex("valid-refresh-token")
		err := repo.StoreRefreshToken(context.Background(), authTestUserID, hash, time.Now().Add(1*time.Hour))
		require.NoError(t, err)

		userID, err := repo.ValidateRefreshToken(context.Background(), hash)
		require.NoError(t, err)
		assert.Equal(t, authTestUserID, userID)
	})

	t.Run("expired token returns empty", func(t *testing.T) {
		hash := sha256Hex("expired-refresh-token")
		err := repo.StoreRefreshToken(context.Background(), authTestUserID, hash, time.Now().Add(-1*time.Hour))
		require.NoError(t, err)

		userID, err := repo.ValidateRefreshToken(context.Background(), hash)
		require.NoError(t, err)
		assert.Empty(t, userID)
	})

	t.Run("non-existent token hash returns empty", func(t *testing.T) {
		userID, err := repo.ValidateRefreshToken(context.Background(), "non-existent-hash")
		require.NoError(t, err)
		assert.Empty(t, userID)
	})
}

func TestAuthRepository_RevokeRefreshToken(t *testing.T) {
	repo, cleanup := setupAuthRepoTest(t)
	defer cleanup()

	t.Run("revoked token returns empty", func(t *testing.T) {
		hash := sha256Hex("revoke-test-token")
		err := repo.StoreRefreshToken(context.Background(), authTestUserID, hash, time.Now().Add(1*time.Hour))
		require.NoError(t, err)

		err = repo.RevokeRefreshToken(context.Background(), hash)
		require.NoError(t, err)

		userID, err := repo.ValidateRefreshToken(context.Background(), hash)
		require.NoError(t, err)
		assert.Empty(t, userID)
	})
}

func TestAuthRepository_RevokeAllUserTokens(t *testing.T) {
	repo, cleanup := setupAuthRepoTest(t)
	defer cleanup()

	t.Run("revoke all tokens for user", func(t *testing.T) {
		hash1 := sha256Hex("revoke-all-1")
		hash2 := sha256Hex("revoke-all-2")
		_ = repo.StoreRefreshToken(context.Background(), authTestUserID, hash1, time.Now().Add(1*time.Hour))
		_ = repo.StoreRefreshToken(context.Background(), authTestUserID, hash2, time.Now().Add(1*time.Hour))

		err := repo.RevokeAllUserTokens(context.Background(), authTestUserID)
		require.NoError(t, err)

		userID1, _ := repo.ValidateRefreshToken(context.Background(), hash1)
		userID2, _ := repo.ValidateRefreshToken(context.Background(), hash2)
		assert.Empty(t, userID1)
		assert.Empty(t, userID2)
	})
}

func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", h)
}
