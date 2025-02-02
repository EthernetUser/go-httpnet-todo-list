package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HttpServer HttpServerConfig
	Postgres   PostgresConfig
}

type HttpServerConfig struct {
	Addr         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
}

type PostgresConfig struct {
	ConnString string
}

func New() *Config {
	if err := godotenv.Load("./configs/.env"); err != nil {
		panic(err)
	}

	postgresUser := getEnv("POSTGRES_USER")
	postgresPassword := getEnv("POSTGRES_PASSWORD")
	postgresHost := getEnv("POSTGRES_HOST")
	postgresPort := getEnv("POSTGRES_PORT")
	postgresDB := getEnv("POSTGRES_DB")

	postgresConnString := "postgres://" + postgresUser + ":" + postgresPassword + "@" + postgresHost + ":" + postgresPort + "/" + postgresDB

	return &Config{
		HttpServer: HttpServerConfig{
			Addr: getEnv(
				"HTTP_SERVER_ADDRESS",
			) + ":" + getEnv(
				"HTTP_SERVER_PORT",
			),
			WriteTimeout: parseTimeDurationFromEnv("HTTP_SERVER_WRITE_TIMEOUT"),
			ReadTimeout:  parseTimeDurationFromEnv("HTTP_SERVER_READ_TIMEOUT"),
			IdleTimeout:  parseTimeDurationFromEnv("HTTP_SERVER_IDLE_TIMEOUT"),
		},
		Postgres: PostgresConfig{
			ConnString: postgresConnString,
		},
	}
}

func parseTimeDurationFromEnv(key string) time.Duration {
	value := getEnv(key)

	parsedValue, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("failed to parse %s: %v", key, err)
	}

	return parsedValue
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	log.Fatalf("%s is not set", key)
	return ""
}
