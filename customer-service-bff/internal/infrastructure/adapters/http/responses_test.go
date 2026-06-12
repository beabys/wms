package httpadapter

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuccessResponseJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"key": "value"}
	successResponseJSON(w, http.StatusOK, data)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.True(t, body["success"].(bool))
	assert.Equal(t, "value", body["data"].(map[string]interface{})["key"])
}

func TestSuccessResponseJSONWithStruct(t *testing.T) {
	w := httptest.NewRecorder()
	type testData struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	data := testData{ID: "123", Name: "Test"}
	successResponseJSON(w, http.StatusCreated, data)

	assert.Equal(t, http.StatusCreated, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.True(t, body["success"].(bool))

	dataObj := body["data"].(map[string]interface{})
	assert.Equal(t, "123", dataObj["id"])
	assert.Equal(t, "Test", dataObj["name"])
}

func TestErrorResponseJSON(t *testing.T) {
	w := httptest.NewRecorder()
	errorResponseJSON(w, http.StatusBadRequest, errors.New("invalid input"))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.False(t, body["success"].(bool))
	assert.Equal(t, "invalid input", body["error"])
}

func TestErrorResponseJSONUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	errorResponseJSON(w, http.StatusUnauthorized, errors.New("bad token"))

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.False(t, body["success"].(bool))
	assert.Equal(t, "bad token", body["error"])
}

func TestResponseWriterContentType(t *testing.T) {
	w := httptest.NewRecorder()
	successResponseJSON(w, http.StatusOK, "hello")
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestResponseWriterStatus(t *testing.T) {
	w := httptest.NewRecorder()
	errorResponseJSON(w, http.StatusInternalServerError, errors.New("server error"))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
