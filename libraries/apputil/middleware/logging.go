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

func CorrelationId(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) {
			correlationId := request.Header.Get("Correlation-Id")
			if correlationId == "" {
				correlationId = uuid.NewString()
			}

			requestContext := logging.WithValue(
				request.Context(),
				logging.SlogContextKey{Name: "correlation_id"},
				correlationId,
			)
			request = request.WithContext(requestContext)

			responseWriter.Header().Set("Correlation-Id", correlationId)
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

			responseWriterWrapper := newResponseWriterWrapper(responseWriter)

			next.ServeHTTP(responseWriterWrapper, request)

			if responseWriterWrapper.statusCode >= 400 {
				logRequest(request, requestBody)
				logResponse(request.Context(), responseWriterWrapper, time.Since(startTime))
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

func logResponse(ctx context.Context, responseWriterWrapper *responseWriterWrapper, duration time.Duration) {
	slog.InfoContext(
		ctx,
		"response",
		slog.String("duration", duration.String()),
		slog.String("body", string(responseWriterWrapper.body)),
		slog.Any("headers", responseWriterWrapper.Header()),
	)
}
