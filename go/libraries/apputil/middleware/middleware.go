package middleware

import (
	"net/http"
)

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
