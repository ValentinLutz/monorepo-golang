package infastructure

import (
	"context"
	"errors"
	"fmt"
	"monorepo/libraries/apputil/logging"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type ServerConfig struct {
	Port            int    `yaml:"port"`
	KeyPath         string `yaml:"key_path"`
	CertificatePath string `yaml:"certificate_path"`
}

type Server struct {
	*http.Server
	logger *zerolog.Logger
	config *ServerConfig
}

func NewServer(logger *zerolog.Logger, config *ServerConfig, handler http.Handler) *Server {
	server := &http.Server{
		Addr:     fmt.Sprintf(":%d", config.Port),
		Handler:  handler,
		ErrorLog: logging.NewLoggerWrapper(logger).ToLogger(),
	}

	return &Server{
		Server: server,
		logger: logger,
		config: config,
	}
}

func (server *Server) Start() {
	server.logger.Info().
		Int("port", server.config.Port).
		Msg("starting server")

	err := server.ListenAndServeTLS(server.config.CertificatePath, server.config.KeyPath)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		server.logger.Fatal().
			Err(err).
			Msg("failed to start server")
	}
}

func (server *Server) Stop() {
	timeout := time.Second * 20
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	server.logger.Info().
		Float64("timeout", timeout.Seconds()).
		Msg("stopping server")
	err := server.Shutdown(ctx)
	if err != nil {
		server.logger.Fatal().
			Err(err).
			Msg("failed to stop server gracefully")
	}
	server.logger.Info().
		Msg("server stopped")

	// stop other connections like message queue
}
