package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
	pgxuuid "github.com/vgarvardt/pgx-google-uuid/v5"
)

/**
  Limit code to database initialization, connection testing, and shutdown.
*/

type Storage struct {
	Pool *pgxpool.Pool
}

func NewStorage(config StorageConfig) (*Storage, error) {

	cnf, err := pgxpool.ParseConfig(config.FormatDSN())
	if err != nil {
		panic(err)
	}

	cnf.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), cnf)
	if err != nil {
		panic(err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}

	return &Storage{
		Pool: pool,
	}, nil
}

func (s *Storage) Shutdown() error {
	s.Pool.Close()
	return nil
}
