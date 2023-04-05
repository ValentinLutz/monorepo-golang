package infastructure

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type DatabaseConfig struct {
	Host               string        `yaml:"host"`
	Port               int           `yaml:"port"`
	Username           string        `yaml:"user"`
	Password           string        `yaml:"password"`
	Database           string        `yaml:"database"`
	MaxIdleConnections int           `yaml:"max_idle_connections"`
	MaxOpenConnections int           `yaml:"max_open_connections"`
	MaxIdleTime        time.Duration `yaml:"max_idle_time"`
	MaxLifetime        time.Duration `yaml:"max_lifetime"`
}

type Database struct {
	logger *zerolog.Logger
	config *DatabaseConfig
}

func NewDatabase(logger *zerolog.Logger, config *DatabaseConfig) *Database {
	return &Database{
		config: config,
		logger: logger,
	}
}

func (database *Database) Connect() *sqlx.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		database.config.Host, database.config.Port, database.config.Username, database.config.Password, database.config.Database,
	)

	db, err := sqlx.Open("pgx", psqlInfo)
	if err != nil {
		database.logger.Fatal().
			Err(err).
			Msg("failed to connect to database")
	}

	db.SetMaxIdleConns(database.config.MaxIdleConnections)
	db.SetMaxOpenConns(database.config.MaxOpenConnections)

	err = db.Ping()
	if err != nil {
		db.Close()
		database.logger.Fatal().
			Err(err).
			Msg("failed to ping database")
	}

	return db
}
