package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/raganrrlaw/server/internal/utils"
)

func init() {
	goose.AddMigrationNoTxContext(Up000010, Down000010)
}

func Up000010(ctx context.Context, db *sql.DB) error {
	pd, err := utils.HashPassword("12345678")
	if err != nil {
		return err
	}
	query := `
		INSERT INTO store (
			store_username,
			store_name,
			store_address,
			store_email,
			store_contact_number,
			password_digest,
			store_web_url
		)
    	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if _, err := db.ExecContext(
		ctx,
		query,
		"sk_stores",
		"S.K. & Sons Stores",
		"No, 016, Colombo 10, Sri Lanka",
		"sk_stores@yahoo.com",
		"0763445567",
		pd,
		"localhost:9081",
	); err != nil {
		return err
	}
	return nil
}

func Down000010(ctx context.Context, db *sql.DB) error {
	query := "DELETE FROM store"
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
