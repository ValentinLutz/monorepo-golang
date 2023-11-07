package main

import (
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
)

var (
	database *sqlx.DB

	databaseHost = os.Getenv("DB_HOST")
	databasePort = os.Getenv("DB_PORT")
	databaseName = os.Getenv("DB_NAME")
	databaseUser = os.Getenv("DB_USER")
	databasePass = os.Getenv("DB_PASS")
)

var logger = slog.New(
	slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		},
	),
)
