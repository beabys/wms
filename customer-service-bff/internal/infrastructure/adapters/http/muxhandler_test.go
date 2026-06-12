package httpadapter

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beabys/wms/customer-service-bff/internal/api/v1"
	"github.com/beabys/wms/customer-service-bff/internal/infrastructure/adapters/http/context"
	"github.com/beabys/wms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// mockV1Server implements v1.ServerInterface for testing.
type mockV1Server struct {
	v1.Unimplemented
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

func TestCorsMiddlewareWildcard(t *testing.T) {
	m := corsMiddleware([]string{"*"})

	handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Origin", "http://example.com")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, "http://example.com", resp.Header.Get("Access-Control-Allow-Origin"))
}

func TestCorsMiddlewareSpecificOrigin(t *testing.T) {
	m := corsMiddleware([]string{"http://allowed.com"})

	handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	server := httptest.NewServer(handler)
	defer server.Close()

	// Allowed origin
	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	req.Header.Set("Origin", "http://allowed.com")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, "http://allowed.com", resp.Header.Get("Access-Control-Allow-Origin"))
	resp.Body.Close()

	// Disallowed origin
	req, _ = http.NewRequest(http.MethodGet, server.URL, nil)
	req.Header.Set("Origin", "http://evil.com")
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	resp.Body.Close()
}

func TestCorsMiddlewareOptions(t *testing.T) {
	m := corsMiddleware([]string{"*"})

	handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	server := httptest.NewServer(handler)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodOptions, server.URL, nil)
	req.Header.Set("Origin", "http://example.com")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()
}

func TestJwtExtractorMiddleware(t *testing.T) {
	handler := jwtExtractorMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Context().Value(httpctx.ContextKeyJWT)
		if token == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte(token.(string)))
	}))

	server := httptest.NewServer(handler)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	req.Header.Set("Authorization", "Bearer my-test-token")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "my-test-token", string(body))
}

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name string
		auth string
		want string
	}{
		{"valid bearer", "Bearer token123", "token123"},
		{"lowercase bearer", "bearer token123", "token123"},
		{"no auth header", "", ""},
		{"wrong prefix", "Basic dXNlcjpwYXNz", ""},
		{"malformed", "Bearer", ""},
		{"extra spaces", "Bearer  token123  ", " token123  "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.auth != "" {
				r.Header.Set("Authorization", tt.auth)
			}
			got := extractBearerToken(r)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCorsMiddlewareNoOrigin(t *testing.T) {
	m := corsMiddleware([]string{"*"})

	handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	server := httptest.NewServer(handler)
	defer server.Close()

	// Request without Origin header
	resp, err := http.Get(server.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
}
