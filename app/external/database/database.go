package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type Database struct {
	logger *zerolog.Logger
}

func NewDatabase(logger *zerolog.Logger) *Database {
	return &Database{logger: logger}
}

func (database *Database) Connect(config *Config) *sqlx.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database,
	)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		database.logger.Fatal().
			Err(err).
			Msg("Failed to connect to database")
	}

	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		database.logger.Fatal().
			Err(err).
			Msg("Failed to ping database")
	}

	return db
}
