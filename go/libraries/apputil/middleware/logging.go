package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"monorepo/libraries/apputil/logging"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func CorrelationId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		correlationId := request.Header.Get("Correlation-Id")
		if correlationId == "" {
			correlationId = uuid.NewString()
		}

		requestContext := context.WithValue(request.Context(), logging.CorrelationIdKey{}, correlationId)
		request = request.WithContext(requestContext)

		logger := zerolog.Ctx(requestContext)
		logger.UpdateContext(func(loggingContext zerolog.Context) zerolog.Context {
			return loggingContext.Str("correlation_id", correlationId)
		})

		responseWriter.Header().Set("Correlation-Id", correlationId)
		next.ServeHTTP(responseWriter, request)
	})
}

func RequestResponseLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		startTime := time.Now()

		requestBody, err := io.ReadAll(request.Body)
		if err != nil {
			log.Error().
				Err(err).
				Msg("error reading request body")
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

func logRequest(context context.Context, request *http.Request, requestBody []byte) {
	logger := zerolog.Ctx(context)

	logEvent := logger.Info().
		Str("method", request.Method).
		Str("path", request.URL.Path).
		Str("query_params", request.URL.Query().Encode())

	if json.Valid(requestBody) {
		logEvent.RawJSON("body", requestBody)
	} else {
		logEvent.Str("body", string(requestBody))
	}

	logEvent.Msg("request")
}

func logResponse(context context.Context, responseWriterContainer *responseWriterContainer, duration time.Duration) {
	logger := zerolog.Ctx(context)

	logEvent := logger.Info().
		Str("duration", duration.String()).
		Int("status", responseWriterContainer.statusCode)

	if json.Valid(responseWriterContainer.body) {
		logEvent.RawJSON("body", responseWriterContainer.body)
	} else {
		logEvent.Str("body", string(responseWriterContainer.body))
	}

	logEvent.Msg("response")
}