package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationNoTxContext(Up000012, Down000012)
}

func Up000012(ctx context.Context, db *sql.DB) error {
	query := `
		CREATE OR REPLACE FUNCTION update_user_shopping_list_update_time()
		RETURNS TRIGGER AS
		$$
		BEGIN
			UPDATE user_shopping_lists SET updated_at = NOW() WHERE user_id = NEW.user_id;
    		RETURN NEW;
		END;
		$$
		LANGUAGE plpgsql;
	`
	if _, err := db.ExecContext(
		ctx,
		query,
	); err != nil {
		return err
	}
	return nil
}

func Down000012(ctx context.Context, db *sql.DB) error {
	query := "DROP FUNCTION IF EXISTS update_user_input_list_update_time;"
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
