package httpresponse

import (
	"encoding/json"
	"github.com/google/uuid"
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

func NewErrorResponse(request *http.Request, errorCode int, error error) ErrorResponse {
	correlationId, ok := request.Context().Value(logging.CorrelationIdKey).(string)
	if !ok {
		correlationId = uuid.NewString()
	}

	var err string
	if error != nil {
		err = error.Error()
	}

	errorResponse := ErrorResponse{
		Code:          errorCode,
		Message:       &err,
		Method:        request.Method,
		Path:          request.RequestURI,
		Timestamp:     time.Now().UTC(),
		CorrelationId: correlationId,
	}

	return errorResponse
}

func Status(responseWriter http.ResponseWriter, statusCode int) {
	responseWriter.WriteHeader(statusCode)
}

func StatusWithBody(responseWriter http.ResponseWriter, statusCode int, body any) {
	bytes, err := json.Marshal(body)
	if err != nil {
		Status(responseWriter, http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	_, err = responseWriter.Write(bytes)
	if err != nil {
		Status(responseWriter, http.StatusInternalServerError)
		return
	}
}

func ErrorWithBody(responseWriter http.ResponseWriter, statusCode int, body any) {
	bytes, err := json.Marshal(body)
	if err != nil {
		Status(responseWriter, http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
	responseWriter.WriteHeader(statusCode)
	_, err = responseWriter.Write(bytes)
	if err != nil {
		Status(responseWriter, http.StatusInternalServerError)
		return
	}
}
