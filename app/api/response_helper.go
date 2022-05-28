package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type JSONWriter interface {
	ToJSON(writer io.Writer) error
}

type ErrorResponse struct {
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Timestamp time.Time `json:"timestamp"`
	Error     int       `json:"error"`
	Message   string    `json:"message"`
}

func (error ErrorResponse) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(error)
}

func StatusOK(responseWriter http.ResponseWriter, request *http.Request, body JSONWriter) {
	Status(responseWriter, request, http.StatusOK, body)
}

func StatusCreated(responseWriter http.ResponseWriter, request *http.Request, body JSONWriter) {
	Status(responseWriter, request, http.StatusCreated, body)
}

func Status(responseWriter http.ResponseWriter, request *http.Request, statusCode int, body JSONWriter) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	if body != nil {
		err := body.ToJSON(responseWriter)
		if err != nil {
			Error(responseWriter, request, http.StatusInternalServerError, 9009, "panic it's over 9000")
		}
	}
}

func Error(responseWriter http.ResponseWriter, request *http.Request, statusCode int, error int, message string) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
	responseWriter.WriteHeader(statusCode)

	errorResponse := ErrorResponse{
		Method:    request.Method,
		Path:      request.RequestURI,
		Timestamp: time.Now().UTC(),
		Error:     error,
		Message:   message,
	}

	_ = errorResponse.ToJSON(responseWriter)
}
