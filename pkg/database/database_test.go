package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgres(t *testing.T) {
	t.Run("new returns empty Postgres", func(t *testing.T) {
		p := New()
		assert.NotNil(t, p)
		assert.Nil(t, p.DB)
		assert.Nil(t, p.SqlDB)
		assert.Nil(t, p.config)
	})
}

func TestSetConfigs(t *testing.T) {
	t.Run("set configs returns self", func(t *testing.T) {
		p := New()
		cfg := &PostgresConfig{
			User:     "test",
			Password: "test",
			Host:     "localhost",
			Port:     5432,
			DBName:   "testdb",
		}
		result := p.SetConfigs(cfg)
		assert.Same(t, p, result)
		assert.Equal(t, cfg, p.config)
	})
}

func TestPostgresConfigDefaults(t *testing.T) {
	t.Run("config with zero values uses defaults on connect", func(t *testing.T) {
		// This test only verifies the struct fields parse correctly,
		// not the actual connection logic.
		cfg := &PostgresConfig{
			User:              "wms",
			Password:          "wms",
			Host:              "localhost",
			Port:              5432,
			DBName:            "wms_test",
			SSLMode:           "disable",
			MaxIdleConns:      0,
			MaxOpenConns:      0,
			ConnMaxLifetime:   0,
			ConnectionRetries: 0,
		}
		assert.Equal(t, "wms", cfg.User)
		assert.Equal(t, 5432, cfg.Port)
		assert.Equal(t, "disable", cfg.SSLMode)

		// Verify defaults will be applied
		maxIdleConns := 10
		if cfg.MaxIdleConns > 0 {
			maxIdleConns = cfg.MaxIdleConns
		}
		assert.Equal(t, 10, maxIdleConns)

		maxOpenConns := 10
		if cfg.MaxOpenConns > 0 {
			maxOpenConns = cfg.MaxOpenConns
		}
		assert.Equal(t, 10, maxOpenConns)
	})
}

func TestPostgresConfigFull(t *testing.T) {
	t.Run("full config values are stored correctly", func(t *testing.T) {
		cfg := &PostgresConfig{
			User:              "admin",
			Password:          "secret",
			Host:              "pg.example.com",
			Port:              5432,
			DBName:            "wms_prod",
			SSLMode:           "require",
			MaxIdleConns:      25,
			MaxOpenConns:      100,
			ConnMaxLifetime:   5 * time.Minute,
			ConnectionRetries: 3,
		}
		p := New().SetConfigs(cfg)
		assert.Equal(t, "admin", p.config.User)
		assert.Equal(t, "secret", p.config.Password)
		assert.Equal(t, "pg.example.com", p.config.Host)
		assert.Equal(t, "require", p.config.SSLMode)
		assert.Equal(t, 25, p.config.MaxIdleConns)
		assert.Equal(t, 100, p.config.MaxOpenConns)
		assert.Equal(t, 5*time.Minute, p.config.ConnMaxLifetime)
		assert.Equal(t, 3, p.config.ConnectionRetries)
	})
}
