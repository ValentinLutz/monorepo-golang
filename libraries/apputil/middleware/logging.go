package middleware

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"time"
)

type CorrelationIdKey struct {
}

func CorrelationId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		correlationId := request.Header.Get("Correlation-Id")
		if correlationId == "" {
			correlationId = uuid.NewString()
		}

		requestContext := context.WithValue(request.Context(), CorrelationIdKey{}, correlationId)
		request = request.WithContext(requestContext)

		logger := zerolog.Ctx(requestContext)
		logger.UpdateContext(func(loggingContext zerolog.Context) zerolog.Context {
			return loggingContext.Str("correlation_id", correlationId)
		})

		responseWriter.Header().Set("Correlation-Id", correlationId)
		next.ServeHTTP(responseWriter, request)
	})
}

func RequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		startTime := time.Now()

		requestBody, err := io.ReadAll(request.Body)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error reading request body")
		}

		reader := io.NopCloser(bytes.NewBuffer(requestBody))
		request.Body = reader

		responseWriterContainer := newResponseWriterContainer(responseWriter)

		next.ServeHTTP(responseWriterContainer, request)

		if responseWriterContainer.statusCode >= 400 {
			requestContext := request.Context()
			logRequest(requestContext, request, requestBody)
			logResponse(requestContext, responseWriterContainer, time.Since(startTime))
		}
	})
}

// responseWriterContainer is a wrapper for http.ResponseWriter that allows
// the written HTTP statusCode and the written HTTP body to be captured.
type responseWriterContainer struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func newResponseWriterContainer(responseWriter http.ResponseWriter) *responseWriterContainer {
	return &responseWriterContainer{ResponseWriter: responseWriter}
}

func (rwc *responseWriterContainer) WriteHeader(statusCode int) {
	rwc.statusCode = statusCode
	rwc.ResponseWriter.WriteHeader(statusCode)
}

func (rwc *responseWriterContainer) Write(bytes []byte) (int, error) {
	rwc.body = bytes
	return rwc.ResponseWriter.Write(bytes)
}

func logRequest(context context.Context, request *http.Request, requestBody []byte) {
	logger := zerolog.Ctx(context)
	logger.Info().
		Str("method", request.Method).
		Str("path", request.URL.Path).
		Str("body", string(requestBody)).
		Msg("Request")
}

func logResponse(context context.Context, responseWriterContainer *responseWriterContainer, duration time.Duration) {
	logger := zerolog.Ctx(context)
	logger.Info().
		Str("duration", duration.String()).
		Int("status", responseWriterContainer.statusCode).
		Str("body", string(responseWriterContainer.body)).
		Msg("Response")
}
