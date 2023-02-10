package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/danielboakye/go-xm/config"
)

// NewConnection will retrieve credentials from GCP secrets manager using secrets pkg and create
// a connection to the database. if NewConnection encounters any errors, it will return the error
func NewConnection(config config.Configurations) (*sql.DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=%s",
		config.DBHost,
		config.DBUser,
		config.DBPass,
		config.DBName,
		"10")

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return db, nil
}
