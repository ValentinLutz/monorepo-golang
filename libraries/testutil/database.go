package testutil

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"
)

type Database struct {
	*pgxpool.Pool
}

func NewDatabase(config DatabaseConfig) *Database {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		config.Host, config.Port, config.Name, config.User, config.Password,
	)

	ctx := context.Background()

	db, err := pgxpool.New(ctx, psqlInfo)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	return &Database{
		Pool: db,
	}
}

func (db *Database) MustLoadAndExec(path string) {
	query := MustLoadQuery(path)
	db.MustExec(query)
}

func (db *Database) MustExec(query string) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()

	_, err := db.Exec(ctx, query)
	if err != nil {
		panic(fmt.Errorf("failed to execute query: %w", err))
	}
}

func MustLoadQuery(path string) string {
	query, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %w", err))
	}

	return string(query)
}
