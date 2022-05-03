package main

import (
	"app/api/orders"
	"app/external/database"
	"app/internal"
	"app/serve"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

const ConfigPath = "config.yaml"

func main() {
	logger := internal.NewLogger()

	newConfig, err := internal.NewConfig(ConfigPath)
	if err != nil {
		logger.Fatal().
			Err(err).
			Str("path", ConfigPath).
			Msg("Failed to load config file")
	}

	internal.SetLogLevel(newConfig.Logger.Level)

	newDatabase := database.NewDatabase(&logger)
	db := newDatabase.Connect(&newConfig.Database)

	server := NewServer(&logger, &newConfig, db)

	logger.Info().
		Str("address", server.Addr).
		Msg("Starting server")
	err = server.ListenAndServe()
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("Failed to start server")
	}
}

func NewServer(logger *zerolog.Logger, config *internal.Config, db *sqlx.DB) *http.Server {
	router := httprouter.New()

	orderAPI := orders.NewAPI(logger, db, config)
	orderAPI.RegisterHandlers(router)

	swaggerUI := serve.NewUI(logger)
	swaggerUI.RegisterUI(router)
	swaggerUI.RegisterOpenAPISchemas(router)

	serverConfig := config.Server

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", serverConfig.Port),
		Handler: router,
		//ErrorLog:     logger,
		ReadTimeout:  time.Second * time.Duration(serverConfig.Timeout.Read),
		WriteTimeout: time.Second * time.Duration(serverConfig.Timeout.Write),
		IdleTimeout:  time.Second * time.Duration(serverConfig.Timeout.Idle),
	}
}
