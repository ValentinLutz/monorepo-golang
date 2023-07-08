package middleware

import (
	"net/http"
)

// responseWriterWrapper is a wrapper for http.ResponseWriter that allows
// the written HTTP statusCode and the written HTTP body to be captured.
type responseWriterWrapper struct {
	http.ResponseWriter

	statusCode int
	body       []byte
}

func newResponseWriterWrapper(responseWriter http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{ResponseWriter: responseWriter}
}

func (responseWriterWrapper *responseWriterWrapper) WriteHeader(statusCode int) {
	responseWriterWrapper.statusCode = statusCode
	responseWriterWrapper.ResponseWriter.WriteHeader(statusCode)
}

func (responseWriterWrapper *responseWriterWrapper) Write(bytes []byte) (int, error) {
	responseWriterWrapper.body = bytes
	return responseWriterWrapper.ResponseWriter.Write(bytes)
}
