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
		config: config,
		server: server,
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

	// stop other connections like message queue
}
