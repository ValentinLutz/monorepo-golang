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
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

const ConfigPath = "config.yaml"

func main() {
	newConfig, err := internal.NewConfig(ConfigPath)
	if err != nil {
		log.Fatal().
			Str("path", ConfigPath).
			Msg("Failed to load config file")
	}

	logger := internal.NewLogger(zerolog.InfoLevel, true)

	newDatabase := database.NewDatabase(&logger)
	db := newDatabase.Connect(&newConfig.Database)

	server := NewServer(&logger, &newConfig.Server, db)

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

func NewServer(logger *zerolog.Logger, serverConfig *internal.ServerConfig, db *sqlx.DB) *http.Server {
	router := httprouter.New()

	orderAPI := orders.NewAPI(logger, db)
	orderAPI.RegisterHandlers(router)

	swaggerUI := serve.NewUI(logger)
	swaggerUI.RegisterUI(router)
	swaggerUI.RegisterOpenAPISchemas(router)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", serverConfig.Port),
		Handler: router,
		//ErrorLog:     logger,
		ReadTimeout:  time.Second * time.Duration(serverConfig.Timeout.Read),
		WriteTimeout: time.Second * time.Duration(serverConfig.Timeout.Write),
		IdleTimeout:  time.Second * time.Duration(serverConfig.Timeout.Idle),
	}
}
