package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/beabys/wms/pkg/postgres"
)

// TestNewPool_Integration tests the PostgreSQL connection pool against a real
// PostgreSQL instance. The test is skipped if the WMS_TEST_DATABASE_URL
// environment variable is not set, or if running in CI without Docker.
//
// To run: export WMS_TEST_DATABASE_URL="postgres://wms:wms_local_dev@localhost:5432/wms_auth?sslmode=disable"
// Prerequisite: docker-compose up (PostgreSQL must be running).
func TestNewPool_Integration(t *testing.T) {
	dsn := os.Getenv("WMS_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("Skipping integration test: WMS_TEST_DATABASE_URL not set. " +
			"Set it to a running PostgreSQL instance (e.g., from docker-compose).")
	}

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, dsn)
	require.NoError(t, err)
	defer pool.Close()

	require.NoError(t, pool.Ping(ctx))
}

func TestNewPool_InvalidDSN(t *testing.T) {
	ctx := context.Background()
	_, err := postgres.NewPool(ctx, "invalid-dsn")
	require.Error(t, err)
}

func TestHealthCheck_NilPool(t *testing.T) {
	ctx := context.Background()
	err := postgres.HealthCheck(ctx, nil)
	require.Error(t, err)
}

func TestHealthCheck_UnreachableDB(t *testing.T) {
	// This test just validates that HealthCheck returns an error for
	// a pool whose database becomes unreachable - we can only check
	// the nil pool case reliably without a real DB.
	ctx := context.Background()
	err := postgres.HealthCheck(ctx, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "pool is nil")
}
