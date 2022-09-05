package api

import (
	"app/internal/errors"
	"app/internal/util"
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

// RequestResponseLogger is a middleware handler that does log requests and responses
// when a client or server error occurs
type RequestResponseLogger struct {
	handler http.Handler
	logger  *util.Logger
}

func (rrl *RequestResponseLogger) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	startTime := time.Now()

	requestContext := request.Context()

	correlationId := request.Header.Get("Correlation-ID")
	if correlationId == "" {
		correlationId = uuid.New().String()
	}
	requestContext = context.WithValue(requestContext, util.CorrelationIdKey{}, correlationId)

	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		rrl.logger.WithContext(requestContext).
			Error().
			Err(err).
			Msg("Error reading request body")
		Error(responseWriter, request, http.StatusInternalServerError, errors.Panic, err.Error())
		return
	}
	reader := io.NopCloser(bytes.NewBuffer(requestBody))
	request.Body = reader

	rw := wrapResponseWriter(responseWriter)

	rrl.handler.ServeHTTP(rw, request.WithContext(requestContext))

	if rw.status >= 400 {
		rrl.logRequest(requestContext, request, requestBody)
		rrl.logResponse(requestContext, startTime, rw)
	}
}

func (rrl *RequestResponseLogger) logRequest(context context.Context, request *http.Request, requestBody []byte) {
	loggingEvent := rrl.logger.WithContext(context).
		Info().
		Str("method", request.Method).
		Str("path", request.URL.Path)

	if !json.Valid(requestBody) {
		loggingEvent.Bool("valid_json", false).
			Str("body", string(requestBody))
	} else {
		loggingEvent.Bool("valid_json", true).
			RawJSON("body", requestBody)
	}

	loggingEvent.Msg("Request")
}

func (rrl *RequestResponseLogger) logResponse(context context.Context, startTime time.Time, rw *responseWriter) {
	loggingEvent := rrl.logger.WithContext(context).
		Info().
		Str("duration", time.Since(startTime).String()).
		Int("status", rw.status)

	if !json.Valid(rw.body) {
		loggingEvent.Bool("valid_json", false).
			Str("body", string(rw.body))
	} else {
		loggingEvent.Bool("valid_json", true).
			RawJSON("body", rw.body)
	}

	loggingEvent.Msg("Response")
}

func NewRequestResponseLogger(handlerToWrap http.Handler, logger *util.Logger) *RequestResponseLogger {
	return &RequestResponseLogger{
		handler: handlerToWrap,
		logger:  logger,
	}
}

// responseWriter is a wrapper for http.ResponseWriter that allows
// the written HTTP status code and
// the written HTTP body to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body = b
	return rw.ResponseWriter.Write(b)
}
