package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"monorepo/libraries/apputil/infastructure"
	"monorepo/libraries/apputil/logging"
	"monorepo/services/frontend/app/config"
	"monorepo/services/frontend/app/incoming"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	handler := newHandler(appConfig)
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

func newHandler(config config.Config) http.Handler {
	router := chi.NewRouter()

	api := incoming.New(config)
	api.RegisterRoutes(router)

	logging.LogRoutes(router)

	return router
}
