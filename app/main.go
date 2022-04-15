package main

import (
	"app/api/orders"
	"app/internal/config"
	"app/internal/database"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"time"
)

const ConfigPath = "config.yaml"

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	newConfig, err := config.NewConfig(ConfigPath)
	if err != nil {
		logger.Fatal("Failed to load config file %s", ConfigPath)
	}

	newDatabase := database.NewDatabase(logger)
	_ = newDatabase.Connect(&newConfig.Database)

	server := NewServer(logger, &newConfig.Server)

	logger.Printf("Starting server on %s", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}
}

func NewServer(logger *log.Logger, serverConfig *config.ServerConfig) *http.Server {
	router := httprouter.New()

	orderAPI := orders.NewAPI(logger)
	orderAPI.RegisterHandlers(router)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", serverConfig.Port),
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  time.Second * time.Duration(serverConfig.Timeout.Read),
		WriteTimeout: time.Second * time.Duration(serverConfig.Timeout.Write),
		IdleTimeout:  time.Second * time.Duration(serverConfig.Timeout.Idle),
	}
}
