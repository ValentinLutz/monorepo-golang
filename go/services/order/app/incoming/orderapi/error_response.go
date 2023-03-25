package orderapi

import (
	"monorepo/libraries/apputil/logging"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func NewErrorResponse(r *http.Request, errorCode int, error error) ErrorResponse {
	correlationId := r.Context().Value(logging.CorrelationIdKey{})
	if correlationId == nil {
		correlationId = uuid.NewString()
	}

	var err string
	if error != nil {
		err = error.Error()
	}

	errorResponse := ErrorResponse{
		Code:          errorCode,
		Message:       &err,
		Method:        r.Method,
		Path:          r.RequestURI,
		Timestamp:     time.Now().UTC(),
		CorrelationId: correlationId.(string),
	}

	return errorResponse
}
