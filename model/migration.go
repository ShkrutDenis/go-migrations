package model

import (
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

type Migration struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Batch     int       `db:"batch"`
	CreatedAt time.Time `db:"created_at"`
}

func CreateMigrationsTable(connection *sqlx.DB) bool {
	_, err := connection.Exec(checkTableExitSQL())
	if err != nil {
		connection.MustExec(creteTableSQL())
		return true
	}
	return false
}

func AddMigrationRaw(connection *sqlx.DB, migration string, lastBatch int) {
	connection.MustExec(addRawSQL(), migration, lastBatch)
}

func MigrationExist(connection *sqlx.DB, name string) bool {
	var raw Migration
	err := connection.Get(&raw, getRawSQL(), name)
	if err != nil {
		return false
	}
	return true
}

func GetLastMigration(connection *sqlx.DB) (bool, *Migration) {
	var raw Migration
	err := connection.Get(&raw, getLastRawSQL())
	if err != nil {
		return false, &raw
	}
	return true, &raw
}

func GetLastMigrations(connection *sqlx.DB, lastBatch int) []*Migration {
	var list []*Migration
	err := connection.Select(&list, getLastRawsSQL(), lastBatch)
	if err != nil {
		panic(err)
	}
	return list
}

func RemoveMigrationRaw(connection *sqlx.DB, migration string) {
	connection.MustExec(removeMigrationSQL(), migration)
}

func GetLastBatch(connection *sqlx.DB) int {
	var raw Migration
	_ = connection.Get(&raw, getLastBatchSQL())
	return raw.Batch
}

func RemoveLastBatch(connection *sqlx.DB, lastBatch int) {
	connection.MustExec(removeBatchSQL(), lastBatch)
}

// SQL queries functions
func checkTableExitSQL() string {
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		return "SELECT * FROM migrations LIMIT 1;"
	case "postgres":
		return "SELECT * FROM migrations LIMIT 1;"
	default:
		panic("Not supported DB driver: " + os.Getenv("DB_DRIVER"))
	}
}

func creteTableSQL() string {
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		return "CREATE TABLE migrations (id int auto_increment, name varchar(255) not null, batch int not null, created_at datetime default current_timestamp not null, constraint migrations_pk primary key (id));"
	case "postgres":
		return "CREATE TABLE migrations (id serial PRIMARY KEY, name varchar(255) not null, batch int not null, created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP not null);"
	default:
		panic("Not supported DB driver: " + os.Getenv("DB_DRIVER"))
	}
}

func getRawSQL() string {
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		return "SELECT * FROM migrations WHERE name=?;"
	case "postgres":
		return "SELECT * FROM migrations WHERE name=$1;"
	default:
		panic("Not supported DB driver: " + os.Getenv("DB_DRIVER"))
	}
}

func addRawSQL() string {
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		return "INSERT INTO migrations (name, batch) VALUES (?, ?);"
	case "postgres":
		return "INSERT INTO migrations (name, batch) VALUES ($1, $2);"
	default:
		panic("Not supported DB driver: " + os.Getenv("DB_DRIVER"))
	}
}

func getLastRawSQL() string {
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		return "SELECT * FROM migrations ORDER BY created_at DESC, id DESC LIMIT 1;"
	case "postgres":
		return "SELECT * FROM migrations ORDER BY created_at DESC, id DESC LIMIT 1;"
	default:
		panic("Not supported DB driver: " + os.Getenv("DB_DRIVER"))
	}
}

func getLastBatchSQL() string {
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		return "SELECT * FROM migrations ORDER BY batch DESC LIMIT 1;"
	case "postgres":
		return "SELECT * FROM migrations ORDER BY batch DESC LIMIT 1;"
	default:
		panic("Not supported DB driver: " + os.Getenv("DB_DRIVER"))
	}
}

func getLastRawsSQL() string {
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		return "SELECT * FROM migrations WHERE batch=? ORDER BY created_at DESC, id DESC;"
	case "postgres":
		return "SELECT * FROM migrations WHERE batch=$1 ORDER BY created_at DESC, id DESC;"
	default:
		panic("Not supported DB driver: " + os.Getenv("DB_DRIVER"))
	}
}

func removeBatchSQL() string {
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		return "DELETE FROM migrations WHERE batch=?;"
	case "postgres":
		return "DELETE FROM migrations WHERE batch=$1;"
	default:
		panic("Not supported DB driver: " + os.Getenv("DB_DRIVER"))
	}
}

func removeMigrationSQL() string {
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		return "DELETE FROM migrations WHERE name=$1"
	case "postgres":
		return ""
	default:
		panic("Not supported DB driver: " + os.Getenv("DB_DRIVER"))
	}
}
