package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	v1 "github.com/beabys/wms/customer-service-bff/internal/api/v1"
)

var (
	successTrue  = true
	successFalse = false
)

func errorResponse(err error) v1.Envelope {
	errMsg := err.Error()
	return v1.Envelope{
		Success: &successFalse,
		Error:   &errMsg,
	}
}

func errorResponseJSON(w http.ResponseWriter, statusCode int, err error) {
	responseWriter(w, statusCode, errorResponse(err))
}

func middlewareRespHandler(w http.ResponseWriter, msg string, statusCode int) {
	errFormatted := errorResponse(errors.New(msg))
	responseWriter(w, statusCode, errFormatted)
}

func successResponseJSON(w http.ResponseWriter, data map[string]interface{}) {
	envelope := v1.Envelope{
		Success: &successTrue,
		Data:    &data,
	}
	responseWriter(w, http.StatusOK, envelope)
}

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

// decodeJSON decodes JSON from reader into the target.
func decodeJSON(r io.Reader, target interface{}) error {
	return json.NewDecoder(r).Decode(target)
}
