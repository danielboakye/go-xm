package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {

	ctx := context.Background()
	// Create context that listens for the interrupt signal from the OS.
	// ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	// defer stop()

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

	// srv := &http.Server{
	// 	Addr:    ":8080",
	// 	Handler: serverHTTP,
	// }

	// // Listen for the interrupt signal.
	// <-ctx.Done()

	// // Restore default behavior on the interrupt signal and notify user of shutdown.
	// stop()
	// log.Println("shutting down gracefully, press Ctrl+C again to force")

	// // The context is used to inform the server it has 5 seconds to finish
	// // the request it is currently handling
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := srv.Shutdown(ctx); err != nil {
	// 	log.Fatal("Server forced to shutdown: ", err)
	// }

	// log.Println("Server exiting")

}
