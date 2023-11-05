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
	slog.Info(
		"starting server",
		slog.Int("port", server.config.Port),
	)

	err := server.ListenAndServeTLS(server.config.CertificatePath, server.config.KeyPath)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error(
			"failed to start server",
			slog.Any("err", err),
		)
		os.Exit(1)
	}
}

func (server *Server) Stop() {
	timeout := time.Second * 20
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	slog.Info(
		"stopping server",
		slog.String("timeout", timeout.String()),
	)

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error(
			"failed to stop server gracefully",
			slog.Any("err", err),
		)
		os.Exit(1)
	}
	slog.Info("server stopped")
}
