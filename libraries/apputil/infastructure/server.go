package infastructure

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type ServerConfig struct {
	Port            int    `yaml:"port"`
	KeyPath         string `yaml:"key_path"`
	CertificatePath string `yaml:"certificate_path"`
}

type Server struct {
	*http.Server
	config ServerConfig
}

func NewServer(config ServerConfig, handler http.Handler, logger *log.Logger) *Server {
	server := &http.Server{
		Addr:     fmt.Sprintf(":%d", config.Port),
		Handler:  handler,
		ErrorLog: logger,
	}
	return &Server{
		Server: server,
		config: config,
	}
}

func (server *Server) Start() {
	slog.With("port", server.config.Port).
		Info("starting server")

	err := server.ListenAndServeTLS(server.config.CertificatePath, server.config.KeyPath)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.With("err", err).
			Error("failed to start server")
		os.Exit(1)
	}
}

func (server *Server) Stop() {
	timeout := time.Second * 20
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	slog.With("timeout", timeout.Seconds()).
		Info("stopping server")

	err := server.Shutdown(ctx)
	if err != nil {
		slog.With("err", err).
			Error("failed to stop server gracefully")
		os.Exit(1)
	}
	slog.Info("server stopped")

	// stop other connections like message queue
}
