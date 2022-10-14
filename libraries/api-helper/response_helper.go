package api

import (
	"encoding/json"
	"github.com/ValentinLutz/monrepo/libraries/errors"
	"io"
	"net/http"
	"time"
)

type CorrelationIdKey struct {
}

type ErrorResponse struct {
	Code          int       `json:"code"`
	Message       *string   `json:"message,omitempty"`
	Method        string    `json:"method"`
	Path          string    `json:"path"`
	Timestamp     time.Time `json:"timestamp"`
	CorrelationId string    `json:"correlation_id"`
}

type JSONWriter interface {
	ToJSON(writer io.Writer) error
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

func (error ErrorResponse) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(error)
}

func Error(responseWriter http.ResponseWriter, request *http.Request, statusCode int, error errors.Error, message string) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
	responseWriter.WriteHeader(statusCode)

	errorResponse := ErrorResponse{
		Code:          int(error),
		Message:       &message,
		Method:        request.Method,
		Path:          request.RequestURI,
		Timestamp:     time.Now().UTC(),
		CorrelationId: request.Context().Value(CorrelationIdKey{}).(string),
	}

	_ = errorResponse.ToJSON(responseWriter)
}
