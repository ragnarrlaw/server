package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/raganrrlaw/server/database"
)

func main() {
	cfg, err := database.NewDatabaseConfig()
	if err != nil {
		log.Fatal("Unable to read configurations: ", err)
		os.Exit(1)
	}

	conn, err := pgx.Connect(context.Background(), cfg.FormatDSN())
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer conn.Close(context.Background())

	// Seed data for users
	for i := 1; i <= 10; i++ {
		_, err := conn.Exec(context.Background(), `
            INSERT INTO users (id, username, first_name, last_name, email, password)
            VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)`,
			fmt.Sprintf("user%d", i),
			fmt.Sprintf("First%d", i),
			fmt.Sprintf("Last%d", i),
			fmt.Sprintf("user%d@example.com", i),
			fmt.Sprintf("password%d", i),
		)
		if err != nil {
			log.Fatal("Error inserting user:", err)
		}
	}

	// Seed data for stores
	for i := 1; i <= 10; i++ {
		location := fmt.Sprintf(`{"latitude": %.6f, "longitude": %.6f}`, 12.34+float64(i)/10, 56.78+float64(i)/10)
		_, err := conn.Exec(context.Background(), `
            INSERT INTO stores (id, name, location)
            VALUES (gen_random_uuid(), $1, $2)`,
			fmt.Sprintf("Store%d", i),
			location,
		)
		if err != nil {
			log.Fatal("Error inserting store:", err)
		}
	}

	fmt.Println("Data seeded successfully")
}
