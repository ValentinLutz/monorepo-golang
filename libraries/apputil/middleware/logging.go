package middleware

import (
	"bytes"
	"context"
	"golang.org/x/exp/slog"
	"io"
	"monorepo/libraries/apputil/logging"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CorrelationId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		correlationId := request.Header.Get("Correlation-Id")
		if correlationId == "" {
			correlationId = uuid.NewString()
		}

		requestContext := context.WithValue(request.Context(), logging.CorrelationIdKey{}, correlationId)
		request = request.WithContext(requestContext)

		responseWriter.Header().Set("Correlation-Id", correlationId)
		next.ServeHTTP(responseWriter, request)
	})
}

func RequestResponseLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		startTime := time.Now()

		requestBody, err := io.ReadAll(request.Body)
		if err != nil {
			slog.With("err", err).
				Error("failed to read request body")
		}

		reader := io.NopCloser(bytes.NewBuffer(requestBody))
		request.Body = reader

		responseWriterWrapper := newResponseWriterWrapper(responseWriter)

		next.ServeHTTP(responseWriterWrapper, request)

		if responseWriterWrapper.statusCode >= 400 {
			logRequest(request, requestBody)
			logResponse(request.Context(), responseWriterWrapper, time.Since(startTime))
		}
	})
}

func logRequest(request *http.Request, requestBody []byte) {
	slog.With("method", request.Method).
		With("path", request.URL.Path).
		With("query_params", request.URL.Query().Encode()).
		With("body", requestBody).
		InfoContext(request.Context(), "request")
}

func logResponse(ctx context.Context, responseWriterWrapper *responseWriterWrapper, duration time.Duration) {
	slog.With("duration", duration.String()).
		With("status", responseWriterWrapper.statusCode).
		With("body", responseWriterWrapper.body).
		InfoContext(ctx, "response")
}
