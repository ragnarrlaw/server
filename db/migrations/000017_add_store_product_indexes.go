package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationNoTxContext(Up000017, Down000017)
}

func Up000017(ctx context.Context, db *sql.DB) error {
	queryId := `CREATE INDEX IF NOT EXISTS store_product_index_on_id ON store_product(id)`
	queryStoreId := `CREATE INDEX IF NOT EXISTS store_product_index_on_store_id ON store_product(store_id)`
	queryStoreProductId := `CREATE UNIQUE INDEX IF NOT EXISTS store_product_index_on_s_p_id ON store_product(store_id, product_id)`
	if _, err := db.ExecContext(ctx, queryId); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, queryStoreId); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, queryStoreProductId); err != nil {
		return err
	}
	return nil
}

func Down000017(ctx context.Context, db *sql.DB) error {
	queryId := `DROP INDEX IF EXISTS store_product_index_on_id`
	queryStoreId := `DROP INDEX IF EXISTS store_product_index_on_store_id`
	queryStoreProductId := `DROP INDEX IF EXISTS store_product_index_on_s_p_id`
	if _, err := db.ExecContext(ctx, queryId); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, queryStoreId); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, queryStoreProductId); err != nil {
		return err
	}
	return nil
}
