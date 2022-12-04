package statusapi

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/hellofresh/health-go/v5"
	psql "github.com/hellofresh/health-go/v5/checks/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"monorepo/services/order/app/config"
	"net/http"
	"time"
)

type API struct {
	db     *sqlx.DB
	config *config.Database
	logger *zerolog.Logger
}

func New(db *sqlx.DB, config *config.Database, logger *zerolog.Logger) *API {
	return &API{
		db:     db,
		config: config,
		logger: logger,
	}
}

func (a *API) RegisterRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Get("/api/status/health", a.registerHealthChecks())
		r.Method("GET", "/api/status/metrics", a.registerPrometheusMetrics())
	})
}

func (a *API) registerHealthChecks() http.HandlerFunc {
	healthStatus, err := health.New()
	if err != nil {
		a.logger.Fatal().
			Err(err).
			Msg("Failed to create health container")
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		a.config.Host, a.config.Port, a.config.Username, a.config.Password, a.config.Database,
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
		a.logger.Fatal().
			Err(err).
			Msg("Failed to create postgres health check")
	}

	return healthStatus.HandlerFunc
}

func (a *API) registerPrometheusMetrics() http.Handler {
	return promhttp.Handler()
}
