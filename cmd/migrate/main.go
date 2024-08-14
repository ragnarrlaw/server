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

	// Connect to PostgreSQL as `bob`
	conn, err := pgx.Connect(context.Background(), cfg.FormatDSNConn())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Create the database if it does not exist
	var dbExists bool
	err = conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)", cfg.GetDatabase()).Scan(&dbExists)
	if err != nil {
		log.Fatalf("Failed to check if database exists: %v", err)
	}

	if !dbExists {
		_, err = conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", cfg.GetDatabase()))
		if err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		fmt.Printf("Database %s created successfully\n", cfg.GetDatabase())
	}

	// Connect to the newly created or existing database
	dbConn, err := pgx.Connect(context.Background(), cfg.FormatDSN())
	if err != nil {
		log.Fatalf("Unable to connect to database %s: %v", cfg.GetDatabase(), err)
	}
	defer dbConn.Close(context.Background())

	// Run migration script to create tables
	createTablesSQL := `
    -- Create the users table
	CREATE TABLE IF NOT EXISTS users (
    	id UUID PRIMARY KEY,
    	username VARCHAR(50) UNIQUE NOT NULL,
    	first_name VARCHAR(50) NOT NULL,
    	last_name VARCHAR(50) NOT NULL,
    	email VARCHAR(100) UNIQUE NOT NULL,
    	password VARCHAR(100) NOT NULL
	);

	-- Create the auth_tokens table
	CREATE TABLE IF NOT EXISTS auth_tokens (
    	id UUID PRIMARY KEY,
    	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    	refresh_token VARCHAR(255) NOT NULL
	);

	-- Create the stores table
	CREATE TABLE IF NOT EXISTS stores (
    	id UUID PRIMARY KEY,
    	name VARCHAR(100) NOT NULL,
    	location JSONB NOT NULL
	);

	-- Create the products table
	CREATE TABLE IF NOT EXISTS products (
    	id UUID PRIMARY KEY,
    	store_id UUID REFERENCES stores(id) ON DELETE CASCADE,
    	name VARCHAR(100) NOT NULL,
    	quantity INT NOT NULL,
    	unit VARCHAR(20) NOT NULL,
    	price NUMERIC(10, 2) NOT NULL
	);
    `

	_, err = dbConn.Exec(context.Background(), createTablesSQL)
	if err != nil {
		log.Fatalf("Failed to execute migration script: %v", err)
	}

	fmt.Println("Migration completed successfully")
}
