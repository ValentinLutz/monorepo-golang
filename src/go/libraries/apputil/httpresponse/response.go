package httpresponse

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"monorepo/libraries/apputil/apierrors"
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

func Status(responseWriter http.ResponseWriter, statusCode int) {
	responseWriter.WriteHeader(statusCode)
}

func StatusWithBody(responseWriter http.ResponseWriter, request *http.Request, statusCode int, body JSONWriter) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	err := body.ToJSON(responseWriter)
	if err != nil {
		StatusInternalServerError(responseWriter, request, err.Error())
	}
}

func StatusOK(responseWriter http.ResponseWriter, request *http.Request, body JSONWriter) {
	StatusWithBody(responseWriter, request, http.StatusOK, body)
}

func StatusCreated(responseWriter http.ResponseWriter, request *http.Request, body JSONWriter) {
	StatusWithBody(responseWriter, request, http.StatusCreated, body)
}

func StatusBadRequest(responseWriter http.ResponseWriter, request *http.Request, message string) {
	Error(responseWriter, request, http.StatusBadRequest, apierrors.BadRequest, message)
}

func StatusUnauthorized(responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("WWW-Authenticate", `Basic realm="order-service"`)
	Status(responseWriter, http.StatusUnauthorized)
}

func StatusInternalServerError(responseWriter http.ResponseWriter, request *http.Request, message string) {
	Error(responseWriter, request, http.StatusInternalServerError, apierrors.Panic, message)
}

func (error ErrorResponse) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(error)
}

func Error(responseWriter http.ResponseWriter, request *http.Request, statusCode int, error apierrors.Error, message string) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
	responseWriter.WriteHeader(statusCode)

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
}
