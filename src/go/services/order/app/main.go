package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"monorepo/libraries/apputil/httpresponse"
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
	"strings"
	"syscall"
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

	newDatabase := infastructure.NewDatabase(&logger, &appConfig.Database)
	db := newDatabase.Connect()

	handler := newHandler(logger, appConfig, db)
	server := infastructure.NewServer(&logger, &appConfig.Server, handler)

	go server.Start()

	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, syscall.SIGINT, syscall.SIGTERM)
	logger.Info().Str("signal", (<-stopChannel).String()).Msg("received signal")

	server.Stop()
}

func newHandler(logger zerolog.Logger, config *config.Config, db *sqlx.DB) http.Handler {
	router := chi.NewRouter()

	orderRepository := orderrepo.NewPostgreSQL(db)
	ordersService := service.NewOrder(db, config, &orderRepository)

	authentication := middleware.Authentication{
		Username: "test",
		Password: "test",
	}

	databaseStats := metrics.NewDatabaseStats(db, metrics.DatabaseOpts{
		Namespace: "app",
		Subsystem: "order",
	})
	prometheus.MustRegister(databaseStats)

	responseTimeHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "app",
		Name:      "http_server_request_duration_seconds",
		Help:      "Histogram of response time in seconds.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "route", "code"})
	prometheus.MustRegister(responseTimeHistogram)

	histogram := middleware.Histogram{Histogram: responseTimeHistogram}

	errorHandler := func(rw http.ResponseWriter, r *http.Request, err error) {
		httpresponse.StatusBadRequest(rw, r, err.Error())
	}

	router.Group(func(r chi.Router) {
		r.Use(hlog.NewHandler(logger))
		r.Use(histogram.Prometheus)
		r.Use(middleware.CorrelationId)
		r.Use(middleware.RequestResponseLogging)
		r.Use(authentication.BasicAuth)
		r.Mount("/api", orderapi.New(ordersService, errorHandler))
	})

	router.Group(func(r chi.Router) {
		statusAPI := statusapi.New(&logger, config, db)
		statusAPI.RegisterRoutes(r)

		openAPI := openapi.New()
		openAPI.RegisterRoutes(r)
	})

	logRoutes(logger, router)

	return router
}

func logRoutes(logger zerolog.Logger, router *chi.Mux) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		logger.Info().
			Str("method", method).
			Str("route", route).
			Msg("register")
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		logger.Error().
			Err(err).
			Msg("failed to walk the routing tree")
	}
}
