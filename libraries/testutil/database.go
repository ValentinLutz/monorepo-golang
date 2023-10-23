package testutil

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	*sqlx.DB
}

func NewDatabase(config *Config) *Database {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.Username,
		config.Database.Password,
		config.Database.Database,
	)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	return &Database{db}
}

func LoadAndExec(db *Database, path string) {
	query := LoadQuery(path)
	Exec(db, query)
}

func LoadQuery(path string) string {
	query, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(query)
}

func Exec(db *Database, query string) {
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}
