package middleware

import (
	"monorepo/libraries/apputil/metrics"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
)

type HttpResponseTimeMetric struct {
	*prometheus.HistogramVec
}

func NewHttpResponseTimeHistogramMetric() *HttpResponseTimeMetric {
	responseTimeHistogram := metrics.NewHttpResponseTimeHistogram(
		metrics.HttpResponseTimeOpts{
			Namespace:  "app",
			LabelNames: []string{"method", "route", "code"},
		},
	)

	return &HttpResponseTimeMetric{
		HistogramVec: responseTimeHistogram,
	}
}

func (httpResponseTimeMetric *HttpResponseTimeMetric) ResponseTimes(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) {
			startTime := time.Now()

			responseWriterContainer := newResponseWriterWrapper(responseWriter)

			next.ServeHTTP(responseWriterContainer, request)

			statusCode := strconv.Itoa(responseWriterContainer.statusCode)
			route := getRoutePattern(request)
			duration := time.Since(startTime)
			httpResponseTimeMetric.WithLabelValues(request.Method, route, statusCode).Observe(duration.Seconds())
		},
	)
}

func getRoutePattern(request *http.Request) string {
	routeContext := chi.RouteContext(request.Context())
	routePattern := routeContext.RoutePattern()
	if routePattern == "" {
		return "undefined"
	}
	return routePattern
}
