//go:build integration

package repository

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/beabys/wms/auth-service/internal/application/auth/repository"
	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
	"github.com/beabys/wms/pkg/database"
	"github.com/beabys/wms/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dbmocks "github.com/beabys/wms/auth-service/mocks/database"
)

type inviteRepoTest struct {
	repo   *InviteRepository
	userID string // valid UUID for FK constraint
	db     *sqlx.DB
}

func setupInviteRepoTest(t *testing.T) *inviteRepoTest {
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

	// Clean up before tests
	_, err = db.Exec("DELETE FROM invite_tokens")
	require.NoError(t, err)
	_, err = db.Exec("DELETE FROM users")
	require.NoError(t, err)

	// Create a test user to satisfy FK constraint
	testUserID := uuid.New().String()
	_, err = db.Exec(`INSERT INTO users (id, email, password_hash, name, role, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, true, NOW(), NOW())`,
		testUserID, "invite-test-admin@test.com", "hash", "Invite Test Admin", "admin",
	)
	require.NoError(t, err)

	log := logger.NewSlogLogger(slog.LevelDebug)
	pg := &database.Postgres{DB: db}
	d := dbmocks.NewDatabase(t)
	d.On("Connect").Return(nil)
	d.On("Ping").Return(nil)
	d.On("Close").Return(nil)
	d.On("GetDBImpl").Return(pg)

	repo := NewInviteRepository(log, d)

	t.Cleanup(func() {
		db.Exec("DELETE FROM invite_tokens")
		db.Exec("DELETE FROM users")
		db.Close()
	})

	return &inviteRepoTest{repo: repo, userID: testUserID, db: db}
}

func TestInviteRepository_Create(t *testing.T) {
	env := setupInviteRepoTest(t)

	t.Run("create invite succeeds", func(t *testing.T) {
		invite, err := model.NewInviteToken("newuser@test.com", env.userID, 24*time.Hour)
		require.NoError(t, err)

		err = env.repo.Create(context.Background(), invite)
		require.NoError(t, err)
		assert.NotEmpty(t, invite.ID)
	})
}

func TestInviteRepository_GetByToken(t *testing.T) {
	env := setupInviteRepoTest(t)

	t.Run("get existing invite by token", func(t *testing.T) {
		invite, err := model.NewInviteToken("getbytoken@test.com", env.userID, 24*time.Hour)
		require.NoError(t, err)
		err = env.repo.Create(context.Background(), invite)
		require.NoError(t, err)

		found, err := env.repo.GetByToken(context.Background(), invite.Token)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, invite.Email, found.Email)
		assert.Equal(t, invite.Token, found.Token)
		assert.Equal(t, model.StatusPending, found.Status)
	})

	t.Run("get non-existent token returns nil", func(t *testing.T) {
		found, err := env.repo.GetByToken(context.Background(), "non-existent-token-value")
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestInviteRepository_Cancel(t *testing.T) {
	env := setupInviteRepoTest(t)

	t.Run("cancel pending invite succeeds", func(t *testing.T) {
		invite, err := model.NewInviteToken("cancel@test.com", env.userID, 24*time.Hour)
		require.NoError(t, err)
		err = env.repo.Create(context.Background(), invite)
		require.NoError(t, err)

		err = env.repo.Cancel(context.Background(), invite.Token)
		require.NoError(t, err)

		found, err := env.repo.GetByToken(context.Background(), invite.Token)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, model.StatusCancelled, found.Status)
	})

	t.Run("cancel already cancelled invite returns error", func(t *testing.T) {
		invite, err := model.NewInviteToken("doublecancel@test.com", env.userID, 24*time.Hour)
		require.NoError(t, err)
		err = env.repo.Create(context.Background(), invite)
		require.NoError(t, err)

		err = env.repo.Cancel(context.Background(), invite.Token)
		require.NoError(t, err)

		err = env.repo.Cancel(context.Background(), invite.Token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found or already processed")
	})

	t.Run("cancel non-existent token returns error", func(t *testing.T) {
		err := env.repo.Cancel(context.Background(), "non-existent-token")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found or already processed")
	})
}

func TestInviteRepository_ListWithStatusFilter(t *testing.T) {
	env := setupInviteRepoTest(t)

	t.Run("list invites filtered by status", func(t *testing.T) {
		pendingInv, _ := model.NewInviteToken("pendinglist@test.com", env.userID, 24*time.Hour)
		_ = env.repo.Create(context.Background(), pendingInv)

		usedInv, _ := model.NewInviteToken("usedlist@test.com", env.userID, 24*time.Hour)
		_ = env.repo.Create(context.Background(), usedInv)
		_ = env.repo.Cancel(context.Background(), usedInv.Token)

		pendingInv2, _ := model.NewInviteToken("pendinglist2@test.com", env.userID, 24*time.Hour)
		_ = env.repo.Create(context.Background(), pendingInv2)

		pendingStatus := model.StatusPending
		result, err := env.repo.List(context.Background(), repository.InviteFilter{
			Status:   &pendingStatus,
			Page:     1,
			PageSize: 50,
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result.Invites), 2)
		for _, inv := range result.Invites {
			assert.Equal(t, model.StatusPending, inv.Status)
		}
	})
}

func TestInviteRepository_ExistsPendingByEmail(t *testing.T) {
	env := setupInviteRepoTest(t)

	t.Run("pending invite exists for email", func(t *testing.T) {
		invite, err := model.NewInviteToken("pendingcheck@test.com", env.userID, 24*time.Hour)
		require.NoError(t, err)
		err = env.repo.Create(context.Background(), invite)
		require.NoError(t, err)

		exists, err := env.repo.ExistsPendingByEmail(context.Background(), "pendingcheck@test.com")
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("no pending invite for email", func(t *testing.T) {
		exists, err := env.repo.ExistsPendingByEmail(context.Background(), "unknown@test.com")
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("existing invite but cancelled returns false", func(t *testing.T) {
		invite, err := model.NewInviteToken("cancelledcheck@test.com", env.userID, 24*time.Hour)
		require.NoError(t, err)
		err = env.repo.Create(context.Background(), invite)
		require.NoError(t, err)

		err = env.repo.Cancel(context.Background(), invite.Token)
		require.NoError(t, err)

		exists, err := env.repo.ExistsPendingByEmail(context.Background(), "cancelledcheck@test.com")
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("duplicate pending email rejection", func(t *testing.T) {
		invite1, err := model.NewInviteToken("duplicate@test.com", env.userID, 24*time.Hour)
		require.NoError(t, err)
		err = env.repo.Create(context.Background(), invite1)
		require.NoError(t, err)

		exists, err := env.repo.ExistsPendingByEmail(context.Background(), "duplicate@test.com")
		require.NoError(t, err)
		assert.True(t, exists)
	})
}

func TestInviteRepository_MarkUsed(t *testing.T) {
	env := setupInviteRepoTest(t)

	t.Run("mark pending invite as used", func(t *testing.T) {
		invite, err := model.NewInviteToken("markused@test.com", env.userID, 24*time.Hour)
		require.NoError(t, err)
		err = env.repo.Create(context.Background(), invite)
		require.NoError(t, err)

		err = env.repo.MarkUsed(context.Background(), invite.Token)
		require.NoError(t, err)

		found, err := env.repo.GetByToken(context.Background(), invite.Token)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, model.StatusUsed, found.Status)
	})

	t.Run("mark already used invite returns error", func(t *testing.T) {
		invite, err := model.NewInviteToken("markusedagain@test.com", env.userID, 24*time.Hour)
		require.NoError(t, err)
		err = env.repo.Create(context.Background(), invite)
		require.NoError(t, err)

		err = env.repo.MarkUsed(context.Background(), invite.Token)
		require.NoError(t, err)

		err = env.repo.MarkUsed(context.Background(), invite.Token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found or already used/cancelled")
	})
}
