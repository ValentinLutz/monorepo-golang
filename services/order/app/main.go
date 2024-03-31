package main

import (
	"flag"
	"log/slog"
	"monorepo/libraries/apputil/infastructure"
	"monorepo/libraries/apputil/logging"
	"monorepo/libraries/apputil/metrics"
	"monorepo/libraries/apputil/middleware"
	"monorepo/services/order/app/config"
	"monorepo/services/order/app/core/service"
	"monorepo/services/order/app/incoming/http/openapi"
	"monorepo/services/order/app/incoming/http/orderapi"
	"monorepo/services/order/app/incoming/http/statusapi"
	"monorepo/services/order/app/outgoing/orderrepo"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	configFile = *flag.String("config", "config/config.yaml", "config file")
)

func main() {
	flag.Parse()

	appConfig, err := config.New(configFile)
	if err != nil {
		panic(err)
	}

	slogHandler := logging.NewSlogHandler(appConfig.Logger)
	contextHandler := logging.NewContextHandler(slogHandler)
	logger := logging.NewSlogLogger(contextHandler)
	slog.SetDefault(logger)

	logLogger := logging.NewLogger(contextHandler, appConfig.Logger)

	database := infastructure.NewDatabase(appConfig.Database)

	handler := newHandler(appConfig, database)
	server := infastructure.NewServer(appConfig.Server, handler, logLogger)

	go server.Start()

	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, syscall.SIGINT, syscall.SIGTERM)
	slog.Info(
		"received signal",
		slog.String("signal", (<-stopChannel).String()),
	)

	server.Stop()
}

func newHandler(config config.Config, database *infastructure.Database) http.Handler {
	router := chi.NewRouter()

	orderRepository := orderrepo.NewPostgreSQL(database)
	ordersService := service.NewOrder(config, &orderRepository)

	authentication := middleware.BasicAuth{
		Username: "test",
		Password: "test",
	}

	databaseStats := metrics.NewDatabaseStats(
		database, metrics.DatabaseOpts{
			Namespace: "app",
			Subsystem: "order",
		},
	)
	prometheus.MustRegister(databaseStats)

	responseTimeHistogramMetric := middleware.NewHttpResponseTimeHistogramMetric()
	prometheus.MustRegister(responseTimeHistogramMetric)

	router.Group(
		func(router chi.Router) {
			router.Use(responseTimeHistogramMetric.ResponseTimes) // before logging
			router.Use(middleware.CorrelationId)                  // before logging
			router.Use(middleware.RequestResponseLogging)
			router.Use(chimiddleware.AllowContentType("application/json"))
			router.Use(authentication.BasicAuth)
			router.Use(middleware.Recover) // always last
			router.Mount("/", orderapi.New(ordersService))
		},
	)

	router.Group(
		func(router chi.Router) {
			statusAPI := statusapi.New(config, database)
			statusAPI.RegisterRoutes(router)

			openAPI := openapi.New()
			openAPI.RegisterRoutes(router)
		},
	)

	logging.LogRoutes(router)

	return router
}
