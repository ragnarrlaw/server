package migrations

import (
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
)

const FILTERED_DATA string = "./db/migrations/filtered_product_sample.csv"
const LEN_UNITS_OF_MEASURE int8 = 5
var UNITS_OF_MEASURE [LEN_UNITS_OF_MEASURE]string = [LEN_UNITS_OF_MEASURE]string{"pcs", "kg", "g", "ml", "l"}

func init() {
	goose.AddMigrationContext(Up000011, Down000011)
}

// super inefficient
func Up000011(ctx context.Context, db *sql.Tx) error {
	file, err := os.Open(FILTERED_DATA)
	if err != nil {
		return err
	}

	parser := csv.NewReader(file)
	parser.Comma = '\t'

	categoryQuery := `INSERT INTO category (category) VALUES ($1) ON CONFLICT (category) DO NOTHING`

	categorySelectQuery := `SELECT id FROM category WHERE category = $1`

	query := `INSERT INTO product (name, brand, brand_tags, category_id, labels, product_quantity, image_url, serving_size, unit_of_measure) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

  r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if _, err := parser.Read(); err == io.EOF {
		return errors.New("empty file provided")
	} else if err != nil {
		return err
	} else {
		for {
			record, err := parser.Read()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return err
			}

			category := strings.Split(record[4], ",")[0]

			if _, err := db.ExecContext(ctx, categoryQuery, category); err != nil {
				db.Rollback()
				return err
			}

			var categoryID uuid.UUID
			if err := db.QueryRowContext(ctx, categorySelectQuery, category).Scan(&categoryID); err != nil {
				db.Rollback()
				return err
			}

			// name, brand, brands_tags, category_id, labels, product_quantity, image_url, serving_size
			if _, err := db.ExecContext(
				ctx,
				query,
				record[1],
				record[2],
				record[3],
				categoryID,
				record[5],
				record[6],
				record[9],
				record[10],
				UNITS_OF_MEASURE[r.Intn(int(LEN_UNITS_OF_MEASURE))],
			); err != nil {
				db.Rollback()
				return err
			}
		}
		return nil
	}
}

func Down000011(ctx context.Context, db *sql.Tx) error {
	query1 := "DELETE FROM product"
	if _, err := db.ExecContext(ctx, query1); err != nil {
		return err
	}
	query2 := "DELETE FROM category"
	if _, err := db.ExecContext(ctx, query2); err != nil {
		return err
	}
	return nil
}
