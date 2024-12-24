package migrations

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-faker/faker/v4"
	"github.com/pressly/goose/v3"
	"github.com/ragnarrlaw/server/internal/utils/hash"
)

const STORE_LOCATION_GEO_DATA = "./db/migrations/store_data.json"

type Properties struct {
	Name            string `json:"name"`
	Shop            string `json:"shop"`
	AddrCity        string `json:"addr:city"`
	AddrHouseNumber string `json:"addr:housenumber"`
	AddrStreet      string `json:"addr:street"`
}

type Geometry struct {
	Coordinates []float64 `json:"coordinates"`
}

type StoreLocationData struct {
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

func init() {
	goose.AddMigrationNoTxContext(Up000010, Down000010)
}

func Up000010(ctx context.Context, db *sql.DB) error {

	file, err := os.Open(STORE_LOCATION_GEO_DATA)

	if err != nil {
		log.Fatalf("file open failed: %s\n", err.Error())
	}
	dec := json.NewDecoder(file)
	// reading array open bracket
	if _, err := dec.Token(); err != nil {
		log.Fatalf("error file parsing the first token : %s\n", err.Error())
	}
	stores := []StoreLocationData{}
	for dec.More() {
		var storeData StoreLocationData
		err := dec.Decode(&storeData)
		if err != nil {
			log.Fatalf("failed during decoding store data %s\n", err.Error())
		}
		stores = append(stores, storeData)
	}

	if _, err := dec.Token(); err != nil {
		log.Fatalf("error file parsing the last token : %v\n", err.Error())
	}
	pd, err := hash.HashPassword("12345678")
	if err != nil {
		return err
	}

	query := `
	INSERT INTO store (
		username,
		name,
		address,
		email,
		location_point,
		contact_number,
		password_digest,
		web_url
	)
	VALUES (
		$3, 
		$4, 
		$5, 
		$6, 
		ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, 
		$7, 
		$8, 
		$9
	)
	`

	for _, store := range stores {
		username := faker.Username()
		if store.Properties.Name == "" {
			store.Properties.Name = fmt.Sprintf("%s %s", username, store.Properties.Shop)
		}

		realAddr := faker.GetRealAddress()
		if store.Properties.AddrCity == "" {
			store.Properties.AddrCity = realAddr.State
		}

		if store.Properties.AddrStreet == "" {
			store.Properties.AddrStreet = realAddr.City
		}

		if _, err := db.ExecContext(
			ctx,
			query,
			store.Geometry.Coordinates[0], // longitude
			store.Geometry.Coordinates[1], // latitude
			username,
			store.Properties.Name,
			fmt.Sprintf("%s, %s", store.Properties.AddrStreet, store.Properties.AddrCity),
			fmt.Sprintf("%s@%s", username, faker.DomainName()),
			faker.E164PhoneNumber(),
			pd,
			"",
		); err != nil {
			return err
		}
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
