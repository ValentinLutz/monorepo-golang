package httpresponse

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"monorepo/libraries/apputil/errors"
	"monorepo/libraries/apputil/logging"
	"net/http"
	"time"
)

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

func StatusUnauthorized(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("WWW-Authenticate", `Basic realm="monke"`)
	Status(responseWriter, request, http.StatusUnauthorized, nil)
}

func Status(responseWriter http.ResponseWriter, request *http.Request, statusCode int, body JSONWriter) {
	if body != nil {
		responseWriter.Header().Set("Content-Type", "application/json")

		err := body.ToJSON(responseWriter)
		if err != nil {
			Error(responseWriter, request, http.StatusInternalServerError, 9009, "panic it's over 9000")
		}
	}
	responseWriter.WriteHeader(statusCode)
}

func (error ErrorResponse) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(error)
}

func InternalServerError(responseWriter http.ResponseWriter, request *http.Request, message string) {
	Error(responseWriter, request, http.StatusInternalServerError, errors.Panic, message)
}

func Error(responseWriter http.ResponseWriter, request *http.Request, statusCode int, error errors.Error, message string) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Header().Set("X-Content-Type-Options", "nosniff")

	correlationId := request.Context().Value(logging.CorrelationIdKey{})
	if correlationId == nil {
		correlationId = uuid.NewString()
	}

	errorResponse := ErrorResponse{
		Code:          int(error),
		Message:       &message,
		Method:        request.Method,
		Path:          request.RequestURI,
		Timestamp:     time.Now().UTC(),
		CorrelationId: correlationId.(string),
	}

	_ = errorResponse.ToJSON(responseWriter)
	responseWriter.WriteHeader(statusCode)
}
