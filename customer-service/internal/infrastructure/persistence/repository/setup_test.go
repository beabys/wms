//go:build integration

package repository

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
	"github.com/beabys/wms/pkg/database"
	"github.com/beabys/wms/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
)

// setupTestDB creates a database connection for integration tests.
func setupTestDB(t *testing.T) database.Database {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnvInt("DB_PORT", 5432)
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "customer_service_test")
	sslMode := getEnv("DB_SSLMODE", "disable")

	db := database.New()
	dbConfig := &database.PostgresConfig{
		User:              dbUser,
		Password:          dbPass,
		Host:              dbHost,
		Port:              dbPort,
		DBName:            dbName,
		SSLMode:           sslMode,
		MaxIdleConns:      2,
		MaxOpenConns:      5,
		ConnMaxLifetime:   5 * time.Minute,
		ConnectionRetries: 3,
	}
	db.SetConfigs(dbConfig)
	if err := db.Connect(); err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	return db
}

// createTestCustomer inserts a test customer and returns its ID.
func createTestCustomer(t *testing.T, db database.Database) string {
	pg, ok := db.GetDBImpl().(*database.Postgres)
	if !ok {
		t.Fatal("expected *database.Postgres")
	}

	id := uuid.New().String()
	now := time.Now()
	email := "test-" + uuid.New().String() + "@example.com"

	query := `INSERT INTO customers (id, company_name, email, status, company_admin_id, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := pg.DB.ExecContext(nil, query,
		id, "Test Customer "+uuid.New().String(), email,
		string(model.CustomerStatusActive), uuid.New().String(), now, now,
	)
	if err != nil {
		t.Fatalf("failed to create test customer: %v", err)
	}
	return id
}

// testLogger returns a no-op logger for tests.
func testLogger() logger.Logger {
	log, _ := logger.NewZapLogger([]string{}, []string{}, zapcore.ErrorLevel)
	return log
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		var val int
		if _, err := fmt.Sscanf(v, "%d", &val); err == nil {
			return val
		}
	}
	return defaultVal
}
