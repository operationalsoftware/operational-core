package db

import (
	"os"
)

// PostgresEnv holds PostgreSQL connection details
type PostgresEnv struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

// LoadPostgresEnv loads PostgreSQL connection details from environment variables
func LoadPostgresEnv() PostgresEnv {
	return PostgresEnv{
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		Host:     os.Getenv("PG_HOST"),
		Port:     os.Getenv("PG_PORT"),
		Database: os.Getenv("PG_DATABASE"),
	}
}
