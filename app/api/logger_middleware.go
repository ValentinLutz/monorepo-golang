package api

import (
	"bytes"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"time"
)

// RequestResponseLogger is a middleware handler that does log requests and responses
// when a client or server error occurs
type RequestResponseLogger struct {
	handler http.Handler
	logger  *zerolog.Logger
}

func (rrl *RequestResponseLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		rrl.logger.Error().Msgf("Error reading request body: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reader := io.NopCloser(bytes.NewBuffer(requestBody))
	r.Body = reader

	rw := wrapResponseWriter(w)

	rrl.handler.ServeHTTP(rw, r)

	if rw.status >= 400 {
		rrl.logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Interface("body", string(requestBody)).
			Msg("Request")
		rrl.logger.Info().
			Str("duration", time.Since(startTime).String()).
			Int("status", rw.status).
			Str("body", string(rw.body)).
			Msg("Response")
	}
}

func NewRequestResponseLogger(handlerToWrap http.Handler, logger *zerolog.Logger) *RequestResponseLogger {
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
