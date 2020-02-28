package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
)

func Connect() (*sqlx.DB, error) {
	return sqlx.Connect(os.Getenv("DB_DRIVER"), getConnectionString())
}

func getConnectionString() string {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_CONNECTION")
	dbName := os.Getenv("DB_NAME")
	return fmt.Sprintf("%v:%v@tcp(%v)/%v?parseTime=true", dbUser, dbPass, dbHost, dbName)
}
