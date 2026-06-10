package httpadapter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/beabys/wms/auth-service-bff/internal/api/v1"
	"github.com/beabys/wms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// mockV1Server implements v1.ServerInterface for testing.
type mockV1Server struct {
	v1.Unimplemented
}

func (m *mockV1Server) InviteUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"success":true,"data":{"token":"test"}}`))
}

func TestNewMuxHandler(t *testing.T) {
	log, err := logger.NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	require.NoError(t, err)

	handler, err := NewMuxHandler(&mockV1Server{}, log, []string{"*"})
	require.NoError(t, err)
	assert.NotNil(t, handler)
}

func TestMuxHandlerHealthEndpoint(t *testing.T) {
	log, _ := logger.NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	handler, err := NewMuxHandler(&mockV1Server{}, log, []string{"*"})
	require.NoError(t, err)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)
	assert.Equal(t, "ok", body["data"].(map[string]interface{})["status"])
	assert.True(t, body["success"].(bool))
}
