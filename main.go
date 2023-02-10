package main

import (
	"log"

	"github.com/danielboakye/go-xm/config"
	"github.com/danielboakye/go-xm/handlers"
	"github.com/danielboakye/go-xm/helpers"
	"github.com/danielboakye/go-xm/pkg/postgres"
	"github.com/danielboakye/go-xm/repo"
	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v4/stdlib"
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

	validator, err := helpers.NewValidation()
	if err != nil {
		log.Fatal("Can't initialize validator")
	}

	db, err := postgres.NewConnection(config)
	if err != nil {
		log.Fatal("Can't connect to the database")
	}
	defer db.Close()

	Repo := repo.NewRepository(db)
	Handler := handlers.NewHandler(Repo, validator, config, nil)

	serverHTTP := NewServerHTTP(Handler, config)
	serverHTTP.Start()
}
