package metrics

import "github.com/prometheus/client_golang/prometheus"

// HttpResponseTimeOpts is a reduced variation of the prometheus.Opts
type HttpResponseTimeOpts struct {
	Namespace  string
	LabelNames []string
}

func NewHttpResponseTimeHistogram(opts HttpResponseTimeOpts) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: opts.Namespace,
			Name:      "http_server_request_duration_seconds",
			Help:      "Histogram of http response times in seconds.",
			Buckets:   prometheus.DefBuckets,
		}, opts.LabelNames,
	)
}
