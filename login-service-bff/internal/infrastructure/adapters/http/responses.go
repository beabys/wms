package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	v1 "github.com/beabys/wms/login-service-bff/internal/api/v1"
)

// errorResponse builds a v1.Error from an error.
func errorResponse(err error) v1.Error {
	errResponse := v1.Error{}
	errResponse.Success = false
	errResponse.Data.Error = err.Error()
	return errResponse
}

// errorResponseJSON writes an error JSON response with the given status code.
func errorResponseJSON(w http.ResponseWriter, statusCode int, err error) {
	responseWriter(w, statusCode, errorResponse(err))
}

// middlewareRespHandler writes a formatted error response.
func middlewareRespHandler(w http.ResponseWriter, msg string, statusCode int) {
	errFormatted := errorResponse(errors.New(msg))
	responseWriter(w, statusCode, errFormatted)
}

// successResponseJSON writes a success JSON response with the given data.
func successResponseJSON(w http.ResponseWriter, data map[string]interface{}) {
	successResponse := v1.Success{}
	successResponse.Success = true
	successResponse.Data = data
	responseWriter(w, http.StatusOK, successResponse)
}

// responseWriter serializes the response and writes it to the http.ResponseWriter.
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
