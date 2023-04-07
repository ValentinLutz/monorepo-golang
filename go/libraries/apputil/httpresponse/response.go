package httpresponse

import (
	"encoding/json"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Code          int       `json:"code"`
	Message       *string   `json:"message,omitempty"`
	Method        string    `json:"method"`
	Path          string    `json:"path"`
	Timestamp     time.Time `json:"timestamp"`
	CorrelationId string    `json:"correlation_id"`
}

func Status(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func StatusWithBody(w http.ResponseWriter, statusCode int, body any) {
	bytes, err := json.Marshal(body)
	if err != nil {
		Status(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(bytes)
	if err != nil {
		Status(w, http.StatusInternalServerError)
		return
	}
}

func ErrorWithBody(w http.ResponseWriter, statusCode int, body any) {
	bytes, err := json.Marshal(body)
	if err != nil {
		Status(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	_, err = w.Write(bytes)
	if err != nil {
		Status(w, http.StatusInternalServerError)
		return
	}
}