package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-faker/faker/v4"
	"github.com/pressly/goose/v3"
	"github.com/raganrrlaw/server/internal/utils/hash"
)

func init() {
	goose.AddMigrationNoTxContext(Up000009, Down000009)
}

func Up000009(ctx context.Context, db *sql.DB) error {
	pd, err := hash.HashPassword("12345678")
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (username, first_name, last_name, email, contact_number, password_digest) 
		VALUES ($1, $2, $3, $4, $5, $6) 
	`
	for i := 0; i < 20; i++ {
		var firstName string
		if i%2 == 0 {
			firstName = faker.FirstNameMale()
		} else {
			firstName = faker.FirstNameFemale()
		}
		j, err := faker.RandomInt(100, 820, 8)
		if err != nil {
			return err
		}
		var lastName string = faker.LastName()
		var contact_number string = faker.E164PhoneNumber()
		var username string = fmt.Sprintf("%s_%s%d", firstName, lastName, j[0])
		var email string = fmt.Sprintf("%s_%s%d@%s", firstName, lastName, j[0], faker.DomainName())
		if _, err := db.ExecContext(
			ctx,
			query,
			username,
			firstName,
			lastName,
			email,
			contact_number,
			pd,
		); err != nil {
			return err
		}
	}
	return nil
}

func Down000009(ctx context.Context, db *sql.DB) error {
	query := "DELETE FROM users"
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
