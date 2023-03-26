package metrics

import "github.com/prometheus/client_golang/prometheus"

func NewResponseTimeHistogram() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "app",
		Name:      "http_server_request_duration_seconds",
		Help:      "Histogram of response time in seconds.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "route", "code"})
}
