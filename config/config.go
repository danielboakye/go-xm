package config

import (
	"os"
	"time"
)

type Configurations struct {
	DBName              string
	DBUser              string
	DBPass              string
	DBHost              string
	AccessTokenDuration time.Duration
	JWTSecretKey        string
	HTTPPort            string
	KafkaURL            string
}

func NewConfigurations() (configs Configurations, err error) {

	at := os.Getenv("ACCESS_TOKEN_DURATION")
	atd, err := time.ParseDuration(at)
	if err != nil {
		return configs, err
	}

	configs = Configurations{
		DBName:              os.Getenv("POSTGRES_DB_NAME"),
		DBUser:              os.Getenv("POSTGRES_DB_USER"),
		DBPass:              os.Getenv("POSTGRES_DB_PASSWORD"),
		DBHost:              os.Getenv("POSTGRES_DB_HOST"),
		AccessTokenDuration: atd,
		JWTSecretKey:        os.Getenv("JWT_SECRET_KEY"),
		HTTPPort:            os.Getenv("HTTP_PORT"),
		KafkaURL:            os.Getenv("KAFKA_URL"),
	}

	return
}
