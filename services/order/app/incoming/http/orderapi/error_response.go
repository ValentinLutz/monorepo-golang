package orderapi

import (
	"monorepo/libraries/apputil/logging"
	"net/http"
	"time"

	"github.com/google/uuid"
)

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
