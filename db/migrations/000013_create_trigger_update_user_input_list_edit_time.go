package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationNoTxContext(Up000013, Down000013)
}

func Up000013(ctx context.Context, db *sql.DB) error {
	query := `
		CREATE OR REPLACE TRIGGER after_user_shopping_list_update
		AFTER UPDATE ON user_shopping_lists
		FOR EACH ROW EXECUTE FUNCTION update_user_shopping_list_update_time();
	`
	if _, err := db.ExecContext(
		ctx,
		query,
	); err != nil {
		return err
	}
	return nil
}

func Down000013(ctx context.Context, db *sql.DB) error {
	query := "DROP TRIGGER IF EXISTS after_user_shopping_list_update ON user_shopping_lists;"
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
