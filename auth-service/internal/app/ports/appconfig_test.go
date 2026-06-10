package ports

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockAppConfig implements AppConfig for testing.
type mockAppConfig struct {
	config map[string]interface{}
}

func (m *mockAppConfig) LoadConfigs() error {
	return nil
}

func (m *mockAppConfig) GetConfigs() interface{} {
	return m.config
}

func TestMockImplementsAppConfig(t *testing.T) {
	var _ AppConfig = (*mockAppConfig)(nil)
}

func TestMockAppConfig(t *testing.T) {
	m := &mockAppConfig{
		config: map[string]interface{}{
			"host": "localhost",
			"port": 8080,
		},
	}

	err := m.LoadConfigs()
	assert.NoError(t, err)

	cfg := m.GetConfigs().(map[string]interface{})
	assert.Equal(t, "localhost", cfg["host"])
	assert.Equal(t, 8080, cfg["port"])
}
