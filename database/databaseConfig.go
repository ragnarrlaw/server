package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	host     string
	port     string
	username string
	password string
	database string
}

func NewDatabaseConfig() (*DatabaseConfig, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	return &DatabaseConfig{
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		username: os.Getenv("DB_USERNAME"),
		password: os.Getenv("DB_PASSWORD"),
		database: os.Getenv("DB_DATABASE"),
	}, nil
}

func (c DatabaseConfig) GetHost() string {
	return c.host
}

func (c DatabaseConfig) GetPort() string {
	return c.port
}

func (c DatabaseConfig) GetUsername() string {
	return c.username
}

func (c DatabaseConfig) GetDatabase() string {
	return c.database
}

func (c DatabaseConfig) FormatDSN() string {
	// "postgres://username:password@localhost:5432/database_name"
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.username, c.password, c.host, c.port, c.database)
}

func (c DatabaseConfig) FormatDSNConn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s", c.username, c.password, c.host, c.port)
}
