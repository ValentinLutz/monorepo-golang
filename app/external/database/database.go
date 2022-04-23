package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type Database struct {
	logger *log.Logger
}

func NewDatabase(logger *log.Logger) *Database {
	return &Database{logger: logger}
}

func (database *Database) Connect(config *Config) *sqlx.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database,
	)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		database.logger.Fatalln(err)
	}

	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		database.logger.Fatalln(err)
	}

	return db
}
