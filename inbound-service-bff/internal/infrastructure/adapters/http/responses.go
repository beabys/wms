package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	v1 "github.com/beabys/wms/inbound-service-bff/internal/api/v1"
)

// errorResponse creates a v1.ApiResponse with an error.
func errorResponse(err error) v1.ApiResponse {
	errMsg := err.Error()
	return v1.ApiResponse{
		Success: boolPtr(false),
		Error:   &errMsg,
	}
}

// errorResponseJSON writes an error response as JSON.
func errorResponseJSON(w http.ResponseWriter, statusCode int, err error) {
	responseWriter(w, statusCode, errorResponse(err))
}

// middlewareRespHandler writes a formatted error response for middleware.
func middlewareRespHandler(w http.ResponseWriter, msg string, statusCode int) {
	errFormatted := errorResponse(errors.New(msg))
	responseWriter(w, statusCode, errFormatted)
}

// successResponseJSON writes a success response as JSON with data.
func successResponseJSON(w http.ResponseWriter, data interface{}) {
	dataMap := map[string]interface{}{}
	jsonData, err := json.Marshal(data)
	if err == nil {
		if err := json.Unmarshal(jsonData, &dataMap); err != nil {
			// fallback: store raw data
			dataMap["value"] = data
		}
	} else {
		dataMap["value"] = data
	}

	resp := v1.ApiResponse{
		Success: boolPtr(true),
		Data:    &dataMap,
	}
	responseWriter(w, http.StatusOK, resp)
}

// responseWriter writes the HTTP response.
func responseWriter(w http.ResponseWriter, statusCode int, response any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	responseData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		responseData = []byte("Something went wrong :(")
	}
	fmt.Fprintln(w, string(responseData))
}

func boolPtr(b bool) *bool {
	return &b
}
