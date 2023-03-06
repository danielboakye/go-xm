package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// config, err := config.NewConfigurations()
	// if err != nil {
	// 	log.Fatal("Can't load configurations")
	// }

	// validator, err := helpers.NewValidation()
	// if err != nil {
	// 	log.Fatal("Can't initialize validator")
	// }

	// db, err := postgres.NewConnection(config)
	// if err != nil {
	// 	log.Fatal("Can't connect to the database")
	// }
	// defer db.Close()

	// r := repo.NewRepository(db)
	// h := handlers.NewHandler(r, validator, config)

	// ihttpHandler := newHTTPHandler(h, config)
	// serverHTTP := newHTTPServer(ihttpHandler, config)

	ctx := context.Background()

	serverHTTP, err := BuildCompileTime(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = serverHTTP.Start()
	if err != nil {
		log.Fatal(err)
	}
}
