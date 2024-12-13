package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationNoTxContext(Up000015, Down000015)
}

func Up000015(ctx context.Context, db *sql.DB) error {
	queryId := `CREATE INDEX IF NOT EXISTS store_index_on_id ON store(id)`
	queryName := `CREATE INDEX IF NOT EXISTS store_index_on_name ON store(name)`
	if _, err := db.ExecContext(ctx, queryId); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, queryName); err != nil {
		return err
	}
	return nil
}

func Down000015(ctx context.Context, db *sql.DB) error {
	queryId := `DROP INDEX IF EXISTS store_index_on_id`
	queryName := `DROP INDEX IF EXISTS store_index_on_name`
	if _, err := db.ExecContext(ctx, queryId); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, queryName); err != nil {
		return err
	}
	return nil
}
