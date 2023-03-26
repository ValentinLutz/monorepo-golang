package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
)

type ResponseTimeMetric struct {
	Histogram *prometheus.HistogramVec
}

func NewResponseTimeMetric(histogram *prometheus.HistogramVec) *ResponseTimeMetric {
	return &ResponseTimeMetric{
		Histogram: histogram,
	}
}

func (metric ResponseTimeMetric) ResponseTimes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		startTime := time.Now()

		responseWriterContainer := newResponseWriterContainer(responseWriter)

		next.ServeHTTP(responseWriterContainer, request)

		duration := time.Since(startTime)
		statusCode := strconv.Itoa(responseWriterContainer.statusCode)
		route := getRoutePattern(request)
		metric.Histogram.WithLabelValues(request.Method, route, statusCode).Observe(duration.Seconds())
	})
}

func getRoutePattern(r *http.Request) string {
	reqContext := chi.RouteContext(r.Context())
	if pattern := reqContext.RoutePattern(); pattern != "" {
		return pattern
	}
	return "undefined"
}
