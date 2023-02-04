package infastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"monorepo/libraries/apputil/logging"
	"monorepo/services/order/app/config"
	"net/http"
	"time"
)

type Server struct {
	logger *zerolog.Logger
	server *http.Server
	config *config.Server
}

func NewServer(logger *zerolog.Logger, serverConfig *config.Server, handler http.Handler) *Server {
	httpServer := &http.Server{
		Addr:     fmt.Sprintf(":%d", serverConfig.Port),
		Handler:  handler,
		ErrorLog: logging.NewLoggerWrapper(logger).ToLogger(),
	}

	return &Server{
		logger: logger,
		server: httpServer,
		config: serverConfig,
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
