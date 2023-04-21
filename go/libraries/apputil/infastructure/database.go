package infastructure

import (
	"fmt"
	"monorepo/libraries/apputil/logging"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	*sqlx.DB
	logger logging.Logger
}

func NewDatabase(logger logging.Logger, config DatabaseConfig) *Database {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database,
	)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("failed to connect to database")
	}

	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetConnMaxIdleTime(config.MaxIdleTime)
	db.SetConnMaxLifetime(config.MaxLifetime)

	err = db.Ping()
	if err != nil {
		db.Close()
		logger.Fatal().
			Err(err).
			Msg("failed to ping database")
	}

	return &Database{
		DB:     db,
		logger: logger,
	}
}
