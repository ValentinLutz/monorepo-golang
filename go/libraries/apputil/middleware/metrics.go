package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"time"
)

type Histogram struct {
	Histogram *prometheus.HistogramVec
}

func (h Histogram) Prometheus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		startTime := time.Now()

		responseWriterContainer := newResponseWriterContainer(responseWriter)

		next.ServeHTTP(responseWriterContainer, request)

		duration := time.Since(startTime)
		statusCode := strconv.Itoa(responseWriterContainer.statusCode)
		route := getRoutePattern(request)
		h.Histogram.WithLabelValues(request.Method, route, statusCode).Observe(duration.Seconds())
	})
}

func getRoutePattern(r *http.Request) string {
	reqContext := chi.RouteContext(r.Context())
	if pattern := reqContext.RoutePattern(); pattern != "" {
		return pattern
	}
	return "undefined"
}
