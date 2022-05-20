package status

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func (api *API) registerPrometheusMetrics() http.Handler {
	return promhttp.Handler()
}
