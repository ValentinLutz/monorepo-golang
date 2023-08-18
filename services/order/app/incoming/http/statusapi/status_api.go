package statusapi

import (
	"fmt"
	"log/slog"
	"monorepo/libraries/apputil/infastructure"
	"monorepo/services/order/app/config"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hellofresh/health-go/v5"
	psql "github.com/hellofresh/health-go/v5/checks/postgres"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type API struct {
	config   config.Config
	database *infastructure.Database
}

func New(config config.Config, database *infastructure.Database) *API {
	return &API{
		config:   config,
		database: database,
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
		slog.With("err", err).
			Error("failed to create health container")
		os.Exit(1)
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
		slog.With("err", err).
			Error("failed to register postgres health check")
		os.Exit(1)
	}

	return healthStatus.HandlerFunc
}

func (api *API) registerPrometheusMetrics() http.Handler {
	return promhttp.Handler()
}
