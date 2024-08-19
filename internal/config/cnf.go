package config

import (
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
}

/*
*

	Load the configurations to environment variables
*/
func Init() error {
	err := godotenv.Load() // load the .env file
	if err != nil {
		log.Fatalf("Error Occurred while loading the env config: %s\n", err.Error())
		return err
	}
	return nil
}

/**
 * Initialize the configuration with default values
 * if the environment variables are not set.

 * @param cfg Config
 * @return error
 * @example cfg := Config{
 *   DBHost: "localhost",
 *   DBPort: "5432",
 *   DBUser: "postgres",
 *   DBPassword: "password",
 *   DBName: "mydb"
 * }
 * err := InitWithDefaults(cfg)
 */

func InitWithDefaults(cfg Config) error {
	v := reflect.ValueOf(cfg)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		envKey := field.Tag.Get("env")
		envValue := v.Field(i).String()

		if os.Getenv(envKey) == "" {
			if err := os.Setenv(envKey, envValue); err != nil {
				log.Fatalf("Error while setting the enirvonment variable %s: %s; error: %s\n", envKey, envValue, err.Error())
				return err
			}
		}
	}
	return nil
}
