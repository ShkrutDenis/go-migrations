package db

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"os"
	"strings"
)

type Config struct {
	user       string
	password   string
	connection string
	host       string
	port       string
	name       string
}

func NewConfig() *Config {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	if host == "" && port == "" {
		connectionParts := strings.Split(os.Getenv("DB_CONNECTION"), ":")
		if len(connectionParts) == 2 {
			host, port = connectionParts[0], connectionParts[1]
		}
	}

	return &Config{
		user:       os.Getenv("DB_USER"),
		password:   os.Getenv("DB_PASSWORD"),
		host:       host,
		port:       port,
		name:       os.Getenv("DB_NAME"),
	}
}