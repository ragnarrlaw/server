package database

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	pgxPool "github.com/jackc/pgx/v5/pgxpool"
)

/*
  Limit code to database operations, initialization, and connection testing.
  And executing the query functions.
*/

type Storage struct {
	pool *pgxPool.Pool
}

func NewStorage(config StorageConfig) (*Storage, error) {
	pool, err := pgxPool.New(context.Background(), config.FormatDSN())
	if err != nil {
		log.Fatalf("Unable to create connection pool: %s\n", err.Error())
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Database connection test failed: %s\n", err.Error())
		return nil, err
	}

	return &Storage{
		pool: pool,
	}, nil
}

func (db *Storage) GetRow(ctx context.Context, sql string, args []interface{}, dest ...interface{}) error {
	row := db.pool.QueryRow(ctx, sql, args...)
	err := row.Scan(dest...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("no rows found")
		}
		return err
	}
	return nil
}

func (db *Storage) GetRows(ctx context.Context, sql string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, col := range rows.FieldDescriptions() {
			rowMap[string(col.Name)] = values[i]
		}
		results = append(results, rowMap)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (db *Storage) Exec(ctx context.Context, sql string, args ...interface{}) error {
	_, err := db.pool.Exec(ctx, sql, args...)
	return err
}

func (s *Storage) Shutdown() error {
	s.pool.Close()
	return nil
}
