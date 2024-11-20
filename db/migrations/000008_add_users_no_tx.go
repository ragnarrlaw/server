package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/raganrrlaw/server/internal/utils"
)

func init() {
	goose.AddMigrationNoTxContext(Up000008, Down000008)
}

func Up000008(ctx context.Context, db *sql.DB) error {
	pd, err := utils.HashPassword("12345678")
	if err != nil {
		return err
	}
	query := `
		INSERT INTO users (username, first_name, last_name, email, contact_number, password_digest) 
		VALUES ($1, $2, $3, $4, $5, $6) 
	`
	if _, err := db.ExecContext(
		ctx,
		query,
		"sapumal_bandara",
		"Sapumal",
		"Bandara",
		"sapumal_banadara@yahoo.com",
		"0771002002",
		pd,
	); err != nil {
		return err
	}
	return nil
}

func Down000008(ctx context.Context, db *sql.DB) error {
	query := "DELETE FROM users"
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
