package main

import (
	"flag"
	"monorepo/libraries/apputil/infastructure"
	"monorepo/libraries/apputil/logging"
	"monorepo/libraries/apputil/metrics"
	"monorepo/libraries/apputil/middleware"
	"monorepo/services/order/app/config"
	"monorepo/services/order/app/core/service"
	"monorepo/services/order/app/incoming/openapi"
	"monorepo/services/order/app/incoming/orderapi"
	"monorepo/services/order/app/incoming/statusapi"
	"monorepo/services/order/app/outgoing/orderrepo"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

var (
	configFile = *flag.String("config", "config/config.yaml", "config file")
)

func main() {
	flag.Parse()

	appConfig, err := config.New(configFile)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("path", configFile).
			Msg("failed to load config file")
	}

	logger := logging.NewLogger(appConfig.Logger)

	database := infastructure.NewDatabase(&logger, &appConfig.Database)

	handler := newHandler(logger, appConfig, database)
	server := infastructure.NewServer(&logger, &appConfig.Server, handler)

	go server.Start()

	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, syscall.SIGINT, syscall.SIGTERM)
	logger.Info().Str("signal", (<-stopChannel).String()).Msg("received signal")

	server.Stop()
}

func newHandler(logger zerolog.Logger, config *config.Config, database *infastructure.Database) http.Handler {
	router := chi.NewRouter()

	orderRepository := orderrepo.NewPostgreSQL(database)
	ordersService := service.NewOrder(config, &orderRepository)

	authentication := middleware.Authentication{
		Username: "test",
		Password: "test",
	}

	databaseStats := metrics.NewDatabaseStats(database, metrics.DatabaseOpts{
		Namespace: "app",
		Subsystem: "order",
	})
	prometheus.MustRegister(databaseStats)

	responseTimeHistogram := metrics.NewResponseTimeHistogram()
	prometheus.MustRegister(responseTimeHistogram)
	responseTimeMetric := middleware.NewResponseTimeMetric(responseTimeHistogram)

	router.Group(func(r chi.Router) {
		r.Use(hlog.NewHandler(logger))
		r.Use(responseTimeMetric.ResponseTimes)
		r.Use(middleware.CorrelationId)
		r.Use(middleware.RequestResponseLogging)
		r.Use(authentication.BasicAuth)
		r.Mount("/", orderapi.New(ordersService))
	})

	router.Group(func(r chi.Router) {
		statusAPI := statusapi.New(&logger, config, database)
		statusAPI.RegisterRoutes(r)

		openAPI := openapi.New()
		openAPI.RegisterRoutes(r)
	})

	logging.LogRoutes(logger, router)

	return router
}
