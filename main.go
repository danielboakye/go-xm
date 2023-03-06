package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {

	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	serverHTTP, err := buildCompileTime(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = serverHTTP.Start()
	if err != nil {
		log.Fatal(err)
	}

}
