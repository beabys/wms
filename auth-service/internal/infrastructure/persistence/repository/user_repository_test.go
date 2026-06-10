//go:build integration

package repository

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/beabys/wms/auth-service/internal/application/auth/repository"
	"github.com/beabys/wms/auth-service/internal/domain/auth/model"
	dbmocks "github.com/beabys/wms/auth-service/mocks/database"
	"github.com/beabys/wms/pkg/database"
	"github.com/beabys/wms/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUserRepoTest(t *testing.T) (*UserRepository, func()) {
	t.Helper()
	// Use testcontainers or env-provided DB
	// For CI, use env vars: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
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
	_, err = db.Exec("DELETE FROM users")
	require.NoError(t, err)

	log := logger.NewSlogLogger(slog.LevelDebug)
	pg := &database.Postgres{DB: db}
	d := dbmocks.NewDatabase(t)
	d.On("Connect").Return(nil)
	d.On("Ping").Return(nil)
	d.On("Close").Return(nil)
	d.On("GetDBImpl").Return(pg)

	repo := NewUserRepository(log, d)
	cleanup := func() {
		db.Exec("DELETE FROM users")
		db.Close()
	}

	return repo, cleanup
}

func TestUserRepository_Create(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	t.Run("create user succeeds", func(t *testing.T) {
		user, err := model.NewUser("create@test.com", "password123", "Create Test", "admin", nil)
		require.NoError(t, err)

		err = repo.Create(context.Background(), user)
		require.NoError(t, err)
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	t.Run("get existing user", func(t *testing.T) {
		user, _ := model.NewUser("getbyid@test.com", "password123", "Get By ID", "admin", nil)
		err := repo.Create(context.Background(), user)
		require.NoError(t, err)

		found, err := repo.GetByID(context.Background(), user.ID)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, user.Email, found.Email)
	})

	t.Run("get non-existent user returns nil", func(t *testing.T) {
		found, err := repo.GetByID(context.Background(), "00000000-0000-0000-0000-000000000000")
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	t.Run("get existing user by email", func(t *testing.T) {
		user, _ := model.NewUser("getbyemail@test.com", "password123", "Get By Email", "admin", nil)
		err := repo.Create(context.Background(), user)
		require.NoError(t, err)

		found, err := repo.GetByEmail(context.Background(), "getbyemail@test.com")
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, user.ID, found.ID)
	})

	t.Run("get non-existent email returns nil", func(t *testing.T) {
		found, err := repo.GetByEmail(context.Background(), "nonexistent@test.com")
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestUserRepository_Update(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	t.Run("update user succeeds", func(t *testing.T) {
		user, _ := model.NewUser("update@test.com", "password123", "Update Test", "admin", nil)
		err := repo.Create(context.Background(), user)
		require.NoError(t, err)

		user.Name = "Updated Name"
		err = repo.Update(context.Background(), user)
		require.NoError(t, err)

		found, _ := repo.GetByID(context.Background(), user.ID)
		require.NotNil(t, found)
		assert.Equal(t, "Updated Name", found.Name)
	})
}

func TestUserRepository_Deactivate(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	t.Run("deactivate user succeeds", func(t *testing.T) {
		user, _ := model.NewUser("deactivate@test.com", "password123", "Deactivate Test", "admin", nil)
		err := repo.Create(context.Background(), user)
		require.NoError(t, err)

		err = repo.Deactivate(context.Background(), user.ID)
		require.NoError(t, err)

		found, _ := repo.GetByID(context.Background(), user.ID)
		require.NotNil(t, found)
		assert.False(t, found.Active)
	})
}

func TestUserRepository_List(t *testing.T) {
	repo, cleanup := setupUserRepoTest(t)
	defer cleanup()

	t.Run("list users with pagination", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			u, _ := model.NewUser(
				"list"+string(rune('0'+i))+"@test.com",
				"password123",
				"List User",
				"admin",
				nil,
			)
			_ = repo.Create(context.Background(), u)
		}

		users, total, err := repo.List(context.Background(), repository.UserFilter{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, 5, total)
		assert.Len(t, users, 5)
	})
}

func getEnvOrDefault(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
