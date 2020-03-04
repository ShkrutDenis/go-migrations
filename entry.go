package go_migrations

import (
	"github.com/ShkrutDenis/go-migrations/db"
	"github.com/ShkrutDenis/go-migrations/model"
	"github.com/ShkrutDenis/go-migrations/store"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
)

var connection *sqlx.DB
var lastBatch int
var firstRun = false

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	connector := db.NewConnector()
	connection, err = connector.Connect()
	if err != nil {
		log.Fatal(err)
	}

	lastBatch = model.GetLastBatch(connection)
	firstRun = model.CreateMigrationsTable(connection)
}

func Migrate(migs []store.Migratable) {
	for _, m := range migs {
		upOrIgnore(m)
	}
}

func Rollback(migs []store.Migratable) {
	store.RegisterMigrations(migs)
	rollBack()
}

func rollBack() {
	forRollback := model.GetLastMigrations(connection, lastBatch)
	for _, m := range forRollback {
		forRun := store.FindMigration(m.Name)
		if forRun == nil {
			log.Fatal("Migration", m.Name, "not found")
		}
		log.Println("Rolling back", forRun.GetName())
		forRun.Down(connection)
		log.Println("Rolled back", forRun.GetName())
	}
	model.RemoveLastBatch(connection, lastBatch)
}

func upOrIgnore(migration store.Migratable) {
	var raw model.Migration
	var err error
	if firstRun {
		goto run
	}
	err = connection.Get(&raw, "SELECT * FROM migrations WHERE name=?", migration.GetName())
	if err != nil {
		goto run
	}
	return
run:
	log.Println("Migrating", migration.GetName())
	migration.Up(connection)
	model.AddMigrationRaw(connection, migration.GetName(), lastBatch+1)
	log.Println("Migrated", migration.GetName())
}
