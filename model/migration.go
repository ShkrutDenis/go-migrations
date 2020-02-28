package model

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type Migration struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Batch     int       `db:"batch"`
	CreatedAt time.Time `db:"created_at"`
}

func CreateMigrationsTable(connection *sqlx.DB) bool {
	_, err := connection.Exec("SELECT * FROM migrations LIMIT 1")
	if err != nil {
		connection.MustExec("CREATE TABLE migrations (id int auto_increment, name varchar(255) not null, batch int not null, created_at datetime default current_timestamp not null, constraint migrations_pk primary key (id));")
		return true
	}
	return false
}

func AddMigrationRaw(connection *sqlx.DB, migration string, lastBatch int) {
	connection.MustExec("INSERT INTO migrations (name, batch) VALUES (?, ?)", migration, lastBatch)
}

func GetLastBatch(connection *sqlx.DB) int {
	var raw Migration
	_ = connection.Get(&raw, "SELECT * FROM migrations ORDER BY batch DESC LIMIT 1")
	return raw.Batch
}

func GetLastMigrations(connection *sqlx.DB, lastBatch int) []*Migration {
	var list []*Migration
	_ = connection.Select(&list, "SELECT * FROM migrations WHERE batch=? ORDER BY created_at DESC;", lastBatch)
	return list
}

func RemoveMigrationRaw(connection *sqlx.DB, migration string) {
	connection.MustExec("DELETE FROM migrations WHERE name=?", migration)
}

func RemoveLastBatch(connection *sqlx.DB, lastBatch int) {
	connection.MustExec("DELETE FROM migrations WHERE batch=?", lastBatch)
}
