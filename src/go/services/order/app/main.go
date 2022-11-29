package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"monorepo/libraries/apputil/logging"
	"monorepo/libraries/apputil/middleware"
	"monorepo/services/order/app/adapter/orderapi"
	"monorepo/services/order/app/adapter/orderitemrepo"
	"monorepo/services/order/app/adapter/orderrepo"
	"monorepo/services/order/app/adapter/statusapi"
	"monorepo/services/order/app/config"
	"monorepo/services/order/app/core/service"
	"monorepo/services/order/app/infastructure"
	"monorepo/services/order/app/serve"
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
	err := server.ListenAndServe()
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
	router := httprouter.New()

	orderRepository := orderrepo.NewPostgreSQL(db)
	orderItemRepository := orderitemrepo.NewPostgreSQL(db)
	ordersService := service.NewOrder(db, config, &orderRepository, &orderItemRepository)

	orderAPI := orderapi.New(config, ordersService)
	orderAPI.RegisterHandlers(router)

	statusAPI := statusapi.New(db, &config.Database)
	statusAPI.RegisterHandlers(router, &logger)

	swaggerUI := serve.NewSwaggerUI()
	swaggerUI.RegisterSwaggerUI(router)
	swaggerUI.RegisterOpenAPISchemas(router)

	handlerChain := alice.New()
	handlerChain = handlerChain.Append(hlog.NewHandler(logger))
	//authentication := middleware.Authentication{
	//	Username: "test",
	//	Password: "test",
	//}
	//handlerChain = handlerChain.Append(authentication.BasicAuth)
	handlerChain = handlerChain.Append(middleware.CorrelationId)
	handlerChain = handlerChain.Append(middleware.RequestLogging)

	handler := handlerChain.Then(router)

	serverConfig := config.Server

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", serverConfig.Port),
		Handler:      handler,
		ErrorLog:     logging.NewLoggerWrapper(&logger).ToLogger(),
		ReadTimeout:  time.Second * time.Duration(serverConfig.Timeout.Read),
		WriteTimeout: time.Second * time.Duration(serverConfig.Timeout.Write),
		IdleTimeout:  time.Second * time.Duration(serverConfig.Timeout.Idle),
	}
}
