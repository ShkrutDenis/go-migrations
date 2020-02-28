package db

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
	"strings"
)

type connectionStringProvider func() string

var connectionProvider map[string]connectionStringProvider = map[string]connectionStringProvider{
	"mysql":    getMysqlConnectionString,
	"postgres": getPostgresConnectionString,
}

func Connect() (*sqlx.DB, error) {
	connectionStringProvider, ok := connectionProvider[os.Getenv("DB_DRIVER")]
	if !ok {
		return nil, errors.New("driver was not provided")
	}

	return sqlx.Connect(os.Getenv("DB_DRIVER"), connectionStringProvider())
}

func getMysqlConnectionString() string {
	dbUser, dbPass, dbHost, dbName := os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_CONNECTION"), os.Getenv("DB_NAME")
	return fmt.Sprintf("%v:%v@tcp(%v)/%v?parseTime=true", dbUser, dbPass, dbHost, dbName)
}

func getPostgresConnectionString() string {
	connectionParts := strings.Split(":", os.Getenv("DB_CONNECTION"))
	if len(connectionParts) < 2 {
		return ""
	}

	dbHost, dbPort := connectionParts[0], connectionParts[1]
	dbUser, dbPass, dbName := os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
}
