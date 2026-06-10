package httpadapter

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func errorResponseJSON(w http.ResponseWriter, statusCode int, err error) {
	response := map[string]interface{}{
		"success": false,
		"error":   err.Error(),
	}
	responseWriter(w, statusCode, response)
}

func successResponseJSON(w http.ResponseWriter, status int, data interface{}) {
	response := map[string]interface{}{
		"success": true,
		"data":    data,
	}
	responseWriter(w, status, response)
}

func responseWriter(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	respData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		respData = []byte("internal server error")
	}
	fmt.Fprintln(w, string(respData))
}
