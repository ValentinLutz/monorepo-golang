package metrics

import "github.com/prometheus/client_golang/prometheus"

// HttpResponseTimeOps is a reduced variation of the prometheus.Opts
type HttpResponseTimeOps struct {
	Namespace  string
	LabelNames []string
}

func NewHttpResponseTimeHistogram(ops HttpResponseTimeOps) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: ops.Namespace,
			Name:      "http_server_request_duration_seconds",
			Help:      "Histogram of http response times in seconds.",
			Buckets:   prometheus.DefBuckets,
		}, ops.LabelNames,
	)
}
