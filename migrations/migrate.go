package main

import (
	"log"
	"os"

	"github.com/danielboakye/go-xm/config"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"

	pgx "github.com/danielboakye/go-xm/pkg/postgres"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config, err := config.NewConfigurations()
	if err != nil {
		log.Fatal("Can't load configurations")
	}

	db, err := pgx.NewConnection(config)
	if err != nil {
		log.Fatal("Can't connect to the database")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migration/sql",
		"postgres", driver)

	if err != nil {
		log.Fatalln(err)
	}

	switch os.Args[1] {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		log.Fatal("Available commands: up, down")
	}

	if err != nil {
		log.Fatalln("Migrate error:", err)
	}
}
