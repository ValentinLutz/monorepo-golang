package infastructure

import (
	"app/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type Database struct {
	config *config.Database
	logger *zerolog.Logger
}

func NewDatabase(config *config.Database, logger *zerolog.Logger) *Database {
	return &Database{config: config, logger: logger}
}

func (database *Database) Connect() *sqlx.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		database.config.Host, database.config.Port, database.config.Username, database.config.Password, database.config.Database,
	)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		database.logger.Fatal().
			Err(err).
			Msg("Failed to connect to database")
	}

	db.SetMaxIdleConns(database.config.MaxIdleConnections)
	db.SetMaxOpenConns(database.config.MaxOpenConnections)

	err = db.Ping()
	if err != nil {
		database.logger.Fatal().
			Err(err).
			Msg("Failed to ping database")
	}

	return db
}
