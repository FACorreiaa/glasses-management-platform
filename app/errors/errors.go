package httperror

import (
	"encoding/json"
	"net/http"
)

type HTTPError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func (e *HTTPError) WriteError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.StatusCode)
	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		return
	}
}

var (
	ErrInvalidID      = &HTTPError{Message: "Invalid ID", StatusCode: http.StatusBadRequest}
	ErrNotFound       = &HTTPError{Message: "Resource not found", StatusCode: http.StatusNotFound}
	ErrInternalServer = &HTTPError{Message: "Internal server error", StatusCode: http.StatusInternalServerError}
	Success           = &HTTPError{Message: "Success", StatusCode: http.StatusOK}
)

func NewHTTPError(message string, statusCode int) *HTTPError {
	return &HTTPError{
		Message:    message,
		StatusCode: statusCode,
	}
}
