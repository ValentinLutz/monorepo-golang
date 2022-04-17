package api

import (
	"net/http"
)

func DefaultResponseHeader(responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Content-Type", "application/json")
}
