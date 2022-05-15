package main

import (
	"app/api/orders"
	"app/api/status"
	"app/external/database"
	"app/internal"
	internalOrders "app/internal/orders"
	"app/serve"
	"context"
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configFile = *flag.String("config", "config/config.yaml", "config file")
	certFile   = *flag.String("cert", "config/cert.crt", "tls certificate file")
	keyFile    = *flag.String("key", "config/cert.key", "tls key file")
)

func main() {
	flag.Parse()
	logger := internal.NewLogger()

	newConfig, err := internal.NewConfig(configFile)
	if err != nil {
		logger.Fatal().
			Err(err).
			Str("path", configFile).
			Msg("Failed to load config file")
	}

	internal.SetLogLevel(newConfig.Logger.Level)

	newDatabase := database.NewDatabase(logger)
	db := newDatabase.Connect(&newConfig.Database)

	server := NewServer(logger, newConfig, db)

	go StartServer(server, logger)
	GracefulServerShutdown(server, logger)
}

func StartServer(server *http.Server, logger *zerolog.Logger) {
	logger.Info().
		Str("address", server.Addr).
		Msg("Starting server")
	err := server.ListenAndServeTLS(certFile, keyFile)
	if err != http.ErrServerClosed {
		logger.Fatal().
			Err(err).
			Msg("Failed to start server")
	}

}

func GracefulServerShutdown(server *http.Server, logger *zerolog.Logger) {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)
	<-osSignal

	timeout := time.Second * 20
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Info().Float64("timeout", timeout.Seconds()).Msg("Stopping server")
	err := server.Shutdown(ctx)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to shutdown server")
	} else {
		logger.Info().Msg("Server stopped")
	}
}

func NewServer(logger *zerolog.Logger, config internal.Config, db *sqlx.DB) *http.Server {
	router := httprouter.New()

	orderRepository := internalOrders.NewOrderRepository(logger, db)
	orderItemRepository := internalOrders.NewOrderItemRepository(logger, db)
	ordersService := internalOrders.NewService(logger, db, config, &orderRepository, &orderItemRepository)

	orderAPI := orders.NewAPI(logger, db, config, ordersService)
	orderAPI.RegisterHandlers(router)

	swaggerUI := serve.NewSwaggerUI(logger)
	swaggerUI.RegisterSwaggerUI(router)
	swaggerUI.RegisterOpenAPISchemas(router)

	statusAPI := status.NewAPI(logger, db, config.Database)
	statusAPI.RegisterHandlers(router)

	serverConfig := config.Server

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", serverConfig.Port),
		Handler:      router,
		ErrorLog:     internal.NewLoggerWrapper(logger).ToLogger(),
		ReadTimeout:  time.Second * time.Duration(serverConfig.Timeout.Read),
		WriteTimeout: time.Second * time.Duration(serverConfig.Timeout.Write),
		IdleTimeout:  time.Second * time.Duration(serverConfig.Timeout.Idle),
	}
}
