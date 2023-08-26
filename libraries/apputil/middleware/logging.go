package middleware

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"monorepo/libraries/apputil/logging"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	CorrelationIdKey = "Correlation-Id"
)

func CorrelationId(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) {
			correlationId := request.Header.Get(CorrelationIdKey)
			if correlationId == "" {
				correlationId = uuid.NewString()
			}

			requestContext := logging.WithValue(
				request.Context(),
				logging.CorrelationIdKey,
				correlationId,
			)
			request = request.WithContext(requestContext)

			responseWriter.Header().Set(CorrelationIdKey, correlationId)
			next.ServeHTTP(responseWriter, request)
		},
	)
}

func RequestResponseLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) {
			startTime := time.Now()

			requestBody, err := io.ReadAll(request.Body)
			if err != nil {
				slog.Error(
					"failed to read request body",
					slog.Any("err", err),
				)
			}
			reader := io.NopCloser(bytes.NewBuffer(requestBody))
			request.Body = reader

			writerWrapper := newResponseWriterWrapper(responseWriter)

			next.ServeHTTP(writerWrapper, request)

			if writerWrapper.statusCode >= 400 {
				logRequest(request, requestBody)
				logResponse(request.Context(), writerWrapper, time.Since(startTime))
			}
		},
	)
}

func logRequest(request *http.Request, requestBody []byte) {
	slog.InfoContext(
		request.Context(),
		"request",
		slog.String("method", request.Method),
		slog.String("path", request.URL.Path),
		slog.String("query_params", request.URL.Query().Encode()),
		slog.String("body", string(requestBody)),
		slog.Any("headers", request.Header),
	)
}

func logResponse(ctx context.Context, responseWriter *responseWriterWrapper, duration time.Duration) {
	slog.InfoContext(
		ctx,
		"response",
		slog.String("duration", duration.String()),
		slog.String("body", string(responseWriter.body)),
		slog.Any("headers", responseWriter.Header()),
	)
}
