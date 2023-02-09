package infastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"monorepo/libraries/apputil/logging"
	"net/http"
	"time"
)

type ServerConfig struct {
	Port            int    `yaml:"port"`
	CertificatePath string `yaml:"certificate_path"`
	KeyPath         string `yaml:"key_path"`
}

type Server struct {
	logger *zerolog.Logger
	config *ServerConfig
	server *http.Server
}

func NewServer(logger *zerolog.Logger, config *ServerConfig, handler http.Handler) *Server {
	server := &http.Server{
		Addr:     fmt.Sprintf(":%d", config.Port),
		Handler:  handler,
		ErrorLog: logging.NewLoggerWrapper(logger).ToLogger(),
	}

	return &Server{
		logger: logger,
		server: server,
		config: config,
	}
}

func (s *Server) Start() {
	s.logger.Info().
		Int("port", s.config.Port).
		Msg("starting server")

	err := s.server.ListenAndServeTLS(s.config.CertificatePath, s.config.KeyPath)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Fatal().
			Err(err).
			Msg("failed to start server")
	}
}

func (s *Server) Stop() {
	timeout := time.Second * 20
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s.logger.Info().
		Float64("timeout", timeout.Seconds()).
		Msg("stopping server")
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Fatal().
			Err(err).
			Msg("failed to stop server gracefully")
	}
	s.logger.Info().
		Msg("server stopped")

	// stop other connections like database, message queue
}