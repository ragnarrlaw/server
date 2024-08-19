package database

import (
	"fmt"
	"os"
)

// CDN Format => "postgres://username:password@localhost:5432/database_name"
// CDN Connection Format => "postgres://username:password@localhost:5432"

type StorageConfig struct {
	host     string
	port     string
	username string
	password string
	database string
}

func NewStorageConfig() (*StorageConfig, error) {
	return &StorageConfig{
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		username: os.Getenv("DB_USERNAME"),
		password: os.Getenv("DB_PASSWORD"),
		database: os.Getenv("DB_DATABASE"),
	}, nil
}

func (c StorageConfig) GetHost() string {
	return c.host
}

func (c StorageConfig) GetPort() string {
	return c.port
}

func (c StorageConfig) GetUsername() string {
	return c.username
}

func (c StorageConfig) GetDatabase() string {
	return c.database
}

func (c StorageConfig) FormatDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.username, c.password, c.host, c.port, c.database)
}

func (c StorageConfig) FormatDSNConn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s", c.username, c.password, c.host, c.port)
}
