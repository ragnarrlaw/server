package database

import (
	"context"
	"log"

	pgxPool "github.com/jackc/pgx/v5/pgxpool"
)

/*
  Limit code to database operations, initialization, and connection testing.
  And executing the query functions.
*/

type DBPool struct {
	dbPool *pgxPool.Pool
}

func NewDBStorage(config DatabaseConfig) (*DBPool, error) {
	pool, err := pgxPool.New(context.Background(), config.FormatDSN())
	if err != nil {
		log.Fatalf("Unable to create connection pool: %s\n", err.Error())
		return nil, err
	}
	return &DBPool{
		dbPool: pool,
	}, nil
}

// initialize the database create the tables
func (d DBPool) Init() error {
	// example:
	query := `
    CREATE TABLE IF NOT EXISTS example_table (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL
    )`
	_, err := d.dbPool.Exec(context.Background(), query)
	if err != nil {
		return err
	}
	return nil
}
