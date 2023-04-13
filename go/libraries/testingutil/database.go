package testingutil

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDatabase(t *testing.T, config *Config) *sqlx.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.Username, config.Database.Password, config.Database.Database,
	)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func LoadAndExec(t *testing.T, db *sqlx.DB, path string) {
	query := LoadQuery(t, path)
	Exec(t, db, query)
}

func LoadQuery(t *testing.T, path string) string {
	query, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return string(query)
}

func Exec(t *testing.T, db *sqlx.DB, query string) {
	_, err := db.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}
