package httpadapter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/beabys/wms/auth-service-bff/mocks/infrastructure/http"
	"github.com/beabys/wms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestNewServer(t *testing.T) {
	mockSvc := mocks.NewAuthHandlerService(t)
	log, _ := logger.NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	s := NewServer(mockSvc, log)
	assert.NotNil(t, s)
}

func TestServerWriteSuccess(t *testing.T) {
	mockSvc := mocks.NewAuthHandlerService(t)
	log, _ := logger.NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	s := NewServer(mockSvc, log)

	w := httptest.NewRecorder()
	data := map[string]string{"msg": "ok"}
	s.writeSuccess(w, data)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.True(t, body["success"].(bool))
}

func TestServerWriteError(t *testing.T) {
	mockSvc := mocks.NewAuthHandlerService(t)
	log, _ := logger.NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	s := NewServer(mockSvc, log)

	w := httptest.NewRecorder()
	s.writeError(w, http.StatusNotFound, "resource not found")

	assert.Equal(t, http.StatusNotFound, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.False(t, body["success"].(bool))
	assert.Equal(t, "resource not found", body["error"])
}
