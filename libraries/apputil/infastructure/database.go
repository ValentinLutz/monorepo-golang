package infastructure

import (
	"fmt"
	"golang.org/x/exp/slog"
	"os"
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
}

func NewDatabase(config DatabaseConfig) *Database {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database,
	)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		slog.With("err", err).
			Error("failed to connect to database")
		os.Exit(1)
	}

	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetConnMaxIdleTime(config.MaxIdleTime)
	db.SetConnMaxLifetime(config.MaxLifetime)

	err = db.Ping()
	if err != nil {
		_ = db.Close()
		slog.With("err", err).Error("failed to ping database")
		os.Exit(1)
	}

	return &Database{
		DB: db,
	}
}
