package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/raganrrlaw/server/internal/utils"
)

func init() {
	goose.AddMigrationNoTxContext(Up000005, Down000005)
}

func Up000005(ctx context.Context, db *sql.DB) error {
	pd, err := utils.HashPassword("12345678")
	if err != nil {
		return err
	}
	query := "INSERT INTO store (store_name, store_address, store_email, store_contact_number, password_digest, store_api_url, store_api_key, store_web_url, store_username) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	// TODO: UDPATE ARGS
	if _, err := db.ExecContext(
		ctx,
		query,
		"sk_stores",
		"", // you can store a geojson object here if needed with the address and all
		"sk_stores@yahoo.com",
		"0763445567",
		pd,
		"localhost:9081",
		"api_key_1",
		"localhost:9081",
		"sk_stores12",
	); err != nil {
		return err
	}
	return nil
}

func Down000005(ctx context.Context, db *sql.DB) error {
	query := "DELETE FROM store"
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
