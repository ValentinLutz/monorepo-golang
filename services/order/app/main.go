package main

import (
	"app/adapter/orderapi"
	"app/adapter/orderitemrepo"
	"app/adapter/orderrepo"
	"app/adapter/statusapi"
	"app/api"
	"app/config"
	"app/core/service"
	"app/infastructure"
	"app/internal/util"
	"app/serve"
	"context"
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configFile = *flag.String("config", "config/config.yaml", "config file")
)

func main() {
	flag.Parse()
	mainLogger := util.New()

	newConfig, err := config.New(configFile)
	if err != nil {
		mainLogger.WithoutContext().
			Fatal().
			Err(err).
			Str("path", configFile).
			Msg("Failed to load config file")
	}

	util.SetLogLevel(newConfig.Logger.Level)

	newDatabase := infastructure.NewDatabase(mainLogger, &newConfig.Database)
	db := newDatabase.Connect()

	server := newServer(mainLogger, newConfig, db)

	go startServer(server, mainLogger)
	shutdownServerGracefully(server, mainLogger)
}

func startServer(server *http.Server, logger *util.Logger) {
	logger.WithoutContext().Info().
		Str("address", server.Addr).
		Msg("Starting server")
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		logger.WithoutContext().Fatal().
			Err(err).
			Msg("Failed to start server")
	}

}

func shutdownServerGracefully(server *http.Server, logger *util.Logger) {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)
	<-osSignal

	timeout := time.Second * 20
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.WithoutContext().
		Info().
		Float64("timeout", timeout.Seconds()).
		Msg("Stopping server")
	err := server.Shutdown(ctx)
	if err != nil {
		logger.WithoutContext().
			Error().
			Err(err).
			Msg("Failed to shutdown server")
	} else {
		logger.WithoutContext().
			Info().
			Msg("Server stopped")
	}
}

func newServer(logger *util.Logger, config *config.Config, db *sqlx.DB) *http.Server {
	router := httprouter.New()

	orderRepository := orderrepo.NewPostgreSQL(logger, db)
	orderItemRepository := orderitemrepo.NewPostgreSQL(logger, db)
	ordersService := service.NewOrder(logger, db, config, &orderRepository, &orderItemRepository)

	orderAPI := orderapi.New(logger, config, ordersService)
	orderAPI.RegisterHandlers(router)

	statusAPI := statusapi.New(logger, db, &config.Database)
	statusAPI.RegisterHandlers(router)

	swaggerUI := serve.NewSwaggerUI(logger)
	swaggerUI.RegisterSwaggerUI(router)
	swaggerUI.RegisterOpenAPISchemas(router)

	// TODO do not use logger middleware on swagger ui and schemas
	routerWithMiddleware := api.NewRequestResponseLogger(router, logger)

	serverConfig := config.Server

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", serverConfig.Port),
		Handler:      routerWithMiddleware,
		ErrorLog:     util.NewLoggerWrapper(logger).ToLogger(),
		ReadTimeout:  time.Second * time.Duration(serverConfig.Timeout.Read),
		WriteTimeout: time.Second * time.Duration(serverConfig.Timeout.Write),
		IdleTimeout:  time.Second * time.Duration(serverConfig.Timeout.Idle),
	}
}
