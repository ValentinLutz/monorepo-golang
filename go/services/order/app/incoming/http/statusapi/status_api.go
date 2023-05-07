package statusapi

import (
	"fmt"
	"monorepo/libraries/apputil/infastructure"
	"monorepo/libraries/apputil/logging"
	"monorepo/services/order/app/config"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hellofresh/health-go/v5"
	psql "github.com/hellofresh/health-go/v5/checks/postgres"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type API struct {
	logger  logging.Logger
	config  config.Config
	databse *infastructure.Database
}

func New(logger logging.Logger, config config.Config, databse *infastructure.Database) *API {
	return &API{
		logger:  logger,
		config:  config,
		databse: databse,
	}
}

func (api *API) RegisterRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Get("/status/health", api.registerHealthChecks())
		r.Method("GET", "/status/metrics", api.registerPrometheusMetrics())
	})
}

func (api *API) registerHealthChecks() http.HandlerFunc {
	healthStatus, err := health.New(health.WithComponent(health.Component{
		Name:    api.config.ServiceName,
		Version: api.config.Version,
	}))
	if err != nil {
		api.logger.Fatal().
			Err(err).
			Msg("failed to create health container")
	}

	databaseConfig := api.config.Database
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		databaseConfig.Host, databaseConfig.Port, databaseConfig.Username, databaseConfig.Password, databaseConfig.Database,
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
		api.logger.Fatal().
			Err(err).
			Msg("failed to create postgres health check")
	}

	return healthStatus.HandlerFunc
}

func (api *API) registerPrometheusMetrics() http.Handler {
	return promhttp.Handler()
}