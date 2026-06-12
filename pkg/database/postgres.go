package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // registers pgx driver via init()
	"github.com/jmoiron/sqlx"
)

// New creates a new Postgres instance.
func New() *Postgres {
	return &Postgres{}
}

// SetConfigs sets the PostgresConfig on the Postgres instance.
func (p *Postgres) SetConfigs(c *PostgresConfig) *Postgres {
	p.config = c
	return p
}

// SetSqlDB sets the underlying *sql.DB.
func (p *Postgres) SetSqlDB(s *sql.DB) *Postgres {
	p.SqlDB = s
	return p
}

// Connect establishes a connection to PostgreSQL using pgx stdlib driver.
func (p *Postgres) Connect() error {
	sslMode := "disable"
	if p.config.SSLMode != "" {
		sslMode = p.config.SSLMode
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.config.User,
		p.config.Password,
		p.config.Host,
		p.config.Port,
		p.config.DBName,
		sslMode,
	)

	// pgx stdlib driver self-registers via init() — use "pgx" directly
	retries := p.config.ConnectionRetries
	if retries <= 0 {
		retries = 3 // default retries
	}
	db, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		for i := 0; i < retries; i++ {
			time.Sleep(time.Second * time.Duration(i+1))
			db, err = sqlx.Connect("pgx", connStr)
			if err != nil {
				continue
			}
			break
		}
	}

	if db == nil {
		return fmt.Errorf("error connecting to postgresql: %w", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	maxIdleConns := 10
	if p.config.MaxIdleConns > 0 {
		maxIdleConns = p.config.MaxIdleConns
	}
	db.DB.SetMaxIdleConns(maxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	maxOpenConns := 10
	if p.config.MaxOpenConns > 0 {
		maxOpenConns = p.config.MaxOpenConns
	}
	db.DB.SetMaxOpenConns(maxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	if p.config.ConnMaxLifetime > 0 {
		db.DB.SetConnMaxLifetime(p.config.ConnMaxLifetime)
	}

	p.SetSqlDB(db.DB)

	if err := p.Ping(); err != nil {
		return err
	}

	p.DB = db

	return nil
}

// Ping verifies a connection to the database is still alive.
func (p *Postgres) Ping() error {
	return p.SqlDB.Ping()
}

// Close closes the database connection.
func (p *Postgres) Close() error {
	if err := p.SqlDB.Close(); err != nil {
		return err
	}
	return nil
}

// GetDBImpl returns the Postgres instance itself.
func (p *Postgres) GetDBImpl() any {
	return p
}
