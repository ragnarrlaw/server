package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationNoTxContext(Up000016, Down000016)
}

func Up000016(ctx context.Context, db *sql.DB) error {
	queryId := `CREATE INDEX IF NOT EXISTS product_index_on_id ON product(id)`
	queryName := `CREATE INDEX IF NOT EXISTS product_index_on_name ON product(name)`
	queryBrand := `CREATE INDEX IF NOT EXISTS product_index_on_brand ON product(brand)`
	if _, err := db.ExecContext(ctx, queryId); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, queryName); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, queryBrand); err != nil {
		return err
	}
	return nil
}

func Down000016(ctx context.Context, db *sql.DB) error {
	queryId := `DROP INDEX IF EXISTS product_index_on_id`
	queryName := `DROP INDEX IF EXISTS product_index_on_name`
	queryBrand := `DROP INDEX IF EXISTS product_index_on_brand`
	if _, err := db.ExecContext(ctx, queryId); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, queryName); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, queryBrand); err != nil {
		return err
	}
	return nil
}
