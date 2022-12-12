package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"monorepo/libraries/apputil/logging"
	"monorepo/libraries/apputil/middleware"
	"monorepo/services/order/app/adapter/openapi"
	"monorepo/services/order/app/adapter/orderapi"
	"monorepo/services/order/app/adapter/orderitemrepo"
	"monorepo/services/order/app/adapter/orderrepo"
	"monorepo/services/order/app/adapter/statusapi"
	"monorepo/services/order/app/config"
	"monorepo/services/order/app/core/service"
	"monorepo/services/order/app/infastructure"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	configFile         = *flag.String("config", "config/config.yaml", "config file")
	tlsCertificateFile = *flag.String("tls_cert", "config/app.crt", "tls certificate file")
	tlsPrivateKeyFile  = *flag.String("tls_key", "config/app.key", "tls private key file")
)

func main() {
	flag.Parse()

	appConfig, err := config.New(configFile)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("path", configFile).
			Msg("Failed to load config file")
	}

	logger := logging.NewLogger(appConfig.Logger)

	newDatabase := infastructure.NewDatabase(&appConfig.Database, &logger)
	db := newDatabase.Connect()

	server := newServer(logger, appConfig, db)

	go startServer(server, &logger)
	shutdownServerGracefully(server, &logger)
}

func startServer(server *http.Server, logger *zerolog.Logger) {
	logger.Info().
		Str("address", server.Addr).
		Msg("Starting server")
	err := server.ListenAndServeTLS(tlsCertificateFile, tlsPrivateKeyFile)
	if err != http.ErrServerClosed {
		logger.Fatal().
			Err(err).
			Msg("Failed to start server")
	}

}

func shutdownServerGracefully(server *http.Server, logger *zerolog.Logger) {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)
	<-osSignal

	timeout := time.Second * 20
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Info().
		Float64("timeout", timeout.Seconds()).
		Msg("Stopping server")
	err := server.Shutdown(ctx)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to shutdown server")
	} else {
		logger.Info().
			Msg("Server stopped")
	}
}

func newServer(logger zerolog.Logger, config *config.Config, db *sqlx.DB) *http.Server {
	router := chi.NewRouter()

	orderRepository := orderrepo.NewPostgreSQL(db)
	orderItemRepository := orderitemrepo.NewPostgreSQL(db)
	ordersService := service.NewOrder(db, config, &orderRepository, &orderItemRepository)

	authentication := middleware.Authentication{
		Username: "test",
		Password: "test",
	}

	responseTimeHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "namespace",
		Name:      "http_server_request_duration_seconds",
		Help:      "Histogram of response time for handler in seconds",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "route", "status_code"})
	prometheus.MustRegister(responseTimeHistogram)

	histogram := middleware.Histogram{Histogram: responseTimeHistogram}

	orderAPI := orderapi.New(config, ordersService)
	router.Group(func(r chi.Router) {
		r.Use(hlog.NewHandler(logger))
		r.Use(middleware.CorrelationId)
		r.Use(authentication.BasicAuth)
		r.Use(middleware.RequestLogging)
		r.Use(histogram.Prometheus)
		r.Mount("/api", orderapi.Handler(orderAPI))

		statusAPI := statusapi.New(db, &config.Database, &logger)
		statusAPI.RegisterRoutes(r)
	})

	router.Group(func(r chi.Router) {
		statusAPI := statusapi.New(db, &config.Database, &logger)
		statusAPI.RegisterRoutes(r)
	})

	openAPI := openapi.New()
	openAPI.RegisterRoutes(router)

	logRoutes(logger, router)

	serverConfig := config.Server
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", serverConfig.Port),
		Handler:      router,
		ErrorLog:     logging.NewLoggerWrapper(&logger).ToLogger(),
		ReadTimeout:  time.Second * time.Duration(serverConfig.Timeout.Read),
		WriteTimeout: time.Second * time.Duration(serverConfig.Timeout.Write),
		IdleTimeout:  time.Second * time.Duration(serverConfig.Timeout.Idle),
	}
}

func logRoutes(logger zerolog.Logger, router *chi.Mux) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		logger.Info().
			Str("method", method).
			Str("route", route).
			Msg("Register")
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to walk")
	}
}
