package main

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func getDatabaseConnection() *sqlx.DB {
	if database != nil {
		return database
	}
	return newDatabase()
}

func newDatabase() *sqlx.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		databaseHost, databasePort, databaseName, databaseUser, databasePass,
	)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		logger.Error(
			"failed to connect to database",
			slog.Any("err", err),
			slog.Group(
				"database",
				slog.String("host", databaseHost),
				slog.String("port", databasePort),
				slog.String("name", databaseName),
			),
		)
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		_ = db.Close()
		logger.Error(
			"failed to ping database",
			slog.Any("err", err),
			slog.Group(
				"database",
				slog.String("host", databaseHost),
				slog.String("port", databasePort),
				slog.String("name", databaseName),
			),
		)
		panic(err)
	}

	return db
}
