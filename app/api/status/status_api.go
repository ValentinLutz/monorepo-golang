package status

import (
	"app/external/database"
	"app/internal/util"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type API struct {
	logger *util.Logger
	db     *sqlx.DB
	config database.Config
}

func NewAPI(logger *util.Logger, db *sqlx.DB, config database.Config) *API {
	return &API{
		logger: logger,
		db:     db,
		config: config,
	}
}

func (api *API) RegisterHandlers(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/api/status/health", api.registerHealthChecks())
	router.Handler(http.MethodGet, "/api/status/metrics", api.registerPrometheusMetrics())
}
