package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // this call registers the pgx driver
	"github.com/pressly/goose/v3"
	"github.com/ragnarrlaw/server/config"
	"github.com/ragnarrlaw/server/db/database"
	_ "github.com/ragnarrlaw/server/db/migrations" // this calls the init functions in the files
	pgxuuid "github.com/vgarvardt/pgx-google-uuid/v5"
)

/**
  Goose doc: https://pressly.github.io/goose/documentation

  To use the sub folder structure to store the migrations in the future with the go embed package
  check this:
      https://stackoverflow.com/questions/66285635/how-do-you-use-go-1-16-embed-features-in-subfolders-packages/67357103#67357103
*/

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
)

func main() {
	log.Println(">>>> running database migrations...")

	flags.Parse(os.Args[1:])

	args := flags.Args()

	flag.Usage = func() {
		log.Println(">>>> usage: go run cmd/migrate/main.go [directory name] [command] [other goose flags]")
		log.Println(">>>> example: go run cmd/migrate/main.go db/migrations up")
	}

	if len(args) < 2 {
		flag.Usage()
		return
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Panic(">>>> error while setting the dialect: ", err.Error())
		return
	}

	if err := config.Init(); err != nil {
		log.Fatalf(">>>> error loading the .env configurations: %v\n", err.Error())
		return
	}

	cn, err := database.NewStorageConfig()
	if err != nil {
		log.Fatalf(">>>> failed creating the storage configuration: %v\n", err.Error())
		return
	}

	cnf, err := pgxpool.ParseConfig(cn.FormatDSN())
	if err != nil {
		log.Fatalf(">>>> error while parsing the configurations: %v\n", err.Error())
		return
	}

	cnf.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	cp := cnf.ConnConfig.ConnString()
	db, err := sql.Open("pgx", cp)
	if err != nil {
		log.Panic(">>>> error failed to open connection: ", err.Error())
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf(">>>> error database connection cannot be closed: %v\n", err.Error())
		}
	}()

	dir, command := args[0], args[1]

	arguments := []string{}
	if len(args) > 2 {
		arguments = append(arguments, args[1:]...)
	}

	switch command {
	case "down-to":
		{
			if len(args) < 3 {
				log.Println(">>>> please provide the version to rollback to")
				return
			} else {
				if version, err := strconv.Atoi(args[2]); err != nil {
					log.Fatalf(">>>> failed to convert the version to integer: %v\n", err.Error())
				} else {
					arguments := []string{strconv.Itoa(version)}
					if err := goose.RunContext(context.Background(), command, db, dir, arguments...); err != nil {
						log.Fatalf(">>>> failed to run the command: %v\n", err.Error())
					}
				}
			}
		}
	default:
		{
			if err := goose.RunContext(context.Background(), command, db, dir, arguments...); err != nil {
				log.Fatalf(">>>> failed to run the command: %v\n", err.Error())
			}
		}
	}
	log.Println(">>>> completed running migrations")
}
