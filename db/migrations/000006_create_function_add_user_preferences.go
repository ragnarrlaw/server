package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationNoTxContext(Up000006, Down000006)
}

func Up000006(ctx context.Context, db *sql.DB) error {
	query := `
		CREATE OR REPLACE FUNCTION add_user_preferences()
		RETURNS TRIGGER AS
		$$
		BEGIN
    		INSERT INTO user_preferences (user_id, preferences)
    		VALUES (NEW.id, '{}');
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

func Down000006(ctx context.Context, db *sql.DB) error {
	query := "DROP FUNCTION IF EXISTS add_user_preferences();"
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
