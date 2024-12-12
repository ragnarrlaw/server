package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/exp/rand"

	"github.com/pressly/goose/v3"
	"github.com/schollz/progressbar/v3"
	"gonum.org/v1/gonum/stat/distuv"
)

type Entity map[string]interface{}

const (
	stockProbability = 0.85
	priceMean        = 10.0
	priceStdDev      = 2.0
)

var (
	stockStatus     = []string{"LIMITED_QUANTITY", "NOT_AVAILABLE", "AVAILABLE"}
	stockStatusProb = []float64{0.1, 0.15, 0.75}
)

func init() {
	goose.AddMigrationContext(Up000014, Down000014)
}

func Up000014(ctx context.Context, db *sql.Tx) error {
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	query := `INSERT INTO store_product (
              store_id,
              product_id,
              price_per_unit,
              listed_unit_of_measure,
              stock_quantity
            ) VALUES ($1, $2, $3, $4, $5)`

	stores, err := readTable("store", ctx, db)
	if err != nil {
		return err
	}

	products, err := readTable("product", ctx, db)
	if err != nil {
		return err
	}

	total := len(stores)
	bar := progressbar.Default(int64(total))

	for _, store := range stores {
		for _, product := range products {
			if rand.Float64() < stockProbability {
				status := sampleStockStatus(r)
				prince := samplePrice(r)
				_, err := db.Exec(
					query,
					store["id"],
					product["id"],
					prince,
					product["unit_of_measure"],
					status,
				)
				if err != nil {
					return err
				}
			}
		}
		bar.Add(1)
	}
	return nil
}

func readTable(tableName string, ctx context.Context, db *sql.Tx) ([]Entity, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	noCols := len(columns)

	var result []Entity
	for rows.Next() {
		e := make([]interface{}, noCols)
		eptr := make([]interface{}, noCols)
		for i := 0; i < noCols; i++ {
			eptr[i] = &e[i]
		}
		if err := rows.Scan(eptr...); err != nil {
			return nil, err
		}
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col] = e[i]
		}
		result = append(result, rowMap)
	}
	return result, nil
}

func sampleStockStatus(src rand.Source) string {
	dist := distuv.NewCategorical(stockStatusProb, src)
	return stockStatus[int64(dist.Rand())]
}

func samplePrice(src rand.Source) float64 {
	dist := distuv.Normal{
		Mu:    priceMean,
		Sigma: priceStdDev,
		Src:   src,
	}
	return dist.Rand()
}

func Down000014(ctx context.Context, db *sql.Tx) error {
	query := `DELETE FROM store_product`
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
