package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationNoTxContext(Up000007, Down000007)
}

func Up000007(ctx context.Context, db *sql.DB) error {
	query := `
		CREATE OR REPLACE TRIGGER after_user_insert
		AFTER INSERT ON users
		FOR EACH ROW
		EXECUTE FUNCTION add_user_preferences();
	`
	if _, err := db.ExecContext(
		ctx,
		query,
	); err != nil {
		return err
	}
	return nil
}

func Down000007(ctx context.Context, db *sql.DB) error {
	query := "DROP TRIGGER IF EXISTS after_user_insert ON users;"
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
