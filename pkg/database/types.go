package database

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Database defines the interface for database operations.
type Database interface {
	Connect() error
	Ping() error
	Close() error
	GetDBImpl() any
}

// Postgres represents a PostgreSQL connection using sqlx + pgx.
type Postgres struct {
	DB    *sqlx.DB
	SqlDB *sql.DB
	config  *PostgresConfig
}

// PostgresConfig holds PostgreSQL connection configuration.
type PostgresConfig struct {
	User            string
	Password        string
	Host            string
	Port            int
	DBName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnectionRetries int
}
