package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockDatabase implements Database interface for testing.
type mockDatabase struct {
	connectErr error
	pingErr    error
	closeErr   error
}

func (m *mockDatabase) Connect() error { return m.connectErr }
func (m *mockDatabase) Ping() error    { return m.pingErr }
func (m *mockDatabase) Close() error   { return m.closeErr }
func (m *mockDatabase) GetDBImpl() any { return m }

func TestMockImplementsDatabase(t *testing.T) {
	var _ Database = (*mockDatabase)(nil)
}

func TestMockDatabaseConnect(t *testing.T) {
	m := &mockDatabase{}
	assert.NoError(t, m.Connect())

	m.connectErr = assert.AnError
	assert.Error(t, m.Connect())
}

func TestMockDatabasePing(t *testing.T) {
	m := &mockDatabase{}
	assert.NoError(t, m.Ping())

	m.pingErr = assert.AnError
	assert.Error(t, m.Ping())
}

func TestMockDatabaseClose(t *testing.T) {
	m := &mockDatabase{}
	assert.NoError(t, m.Close())

	m.closeErr = assert.AnError
	assert.Error(t, m.Close())
}

func TestMockDatabaseGetDBImpl(t *testing.T) {
	m := &mockDatabase{}
	assert.Same(t, m, m.GetDBImpl())
}

func TestPostgresStructInitialization(t *testing.T) {
	p := New()
	assert.NotNil(t, p)
	assert.Nil(t, p.DB)
	assert.Nil(t, p.SqlDB)
	assert.Nil(t, p.config)
}

func TestPostgresConfigInitialization(t *testing.T) {
	cfg := &PostgresConfig{
		User:     "test",
		Password: "test",
		Host:     "localhost",
		Port:     5432,
		DBName:   "testdb",
	}
	assert.Equal(t, "test", cfg.User)
	assert.Equal(t, 5432, cfg.Port)
	assert.Equal(t, "testdb", cfg.DBName)
}
