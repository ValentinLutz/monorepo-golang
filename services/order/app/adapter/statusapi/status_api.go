package statusapi

import (
	"app/config"
	"fmt"
	"github.com/hellofresh/health-go/v4"
	psql "github.com/hellofresh/health-go/v4/checks/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type API struct {
	db     *sqlx.DB
	config *config.Database
}

func New(db *sqlx.DB, config *config.Database) *API {
	return &API{
		db:     db,
		config: config,
	}
}

func (api *API) RegisterHandlers(router *httprouter.Router, logger *zerolog.Logger) {
	router.HandlerFunc(http.MethodGet, "/api/status/health", api.registerHealthChecks(logger))
	router.Handler(http.MethodGet, "/api/status/metrics", api.registerPrometheusMetrics())
}

func (api *API) registerHealthChecks(logger *zerolog.Logger) http.HandlerFunc {
	healthStatus, err := health.New()
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("Failed to create health container")
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		api.config.Host, api.config.Port, api.config.Username, api.config.Password, api.config.Database,
	)
	err = healthStatus.Register(health.Config{
		Name:      "postgresql",
		Timeout:   time.Second * 30,
		SkipOnErr: false,
		Check: psql.New(psql.Config{
			DSN: psqlInfo,
		}),
	})
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("Failed to create postgres health check")
	}

	return healthStatus.HandlerFunc
}

func (api *API) registerPrometheusMetrics() http.Handler {
	return promhttp.Handler()
}
