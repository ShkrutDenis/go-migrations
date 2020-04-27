package go_migrations

import (
	"flag"
	"github.com/ShkrutDenis/go-migrations/config"
	"github.com/ShkrutDenis/go-migrations/db"
	"github.com/ShkrutDenis/go-migrations/model"
	"github.com/ShkrutDenis/go-migrations/store"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
)

var connection *sqlx.DB

func Run(migs []store.Migratable) {
	parseFlags()
	prepare()

	if config.GetConfig().IsRollback {
		store.RegisterMigrations(migs)
		rollBack()
		return
	}

	for _, m := range migs {
		upOrIgnore(m)
	}
}

func rollBack() {
	forRollback := model.GetLastMigrations(connection, config.GetConfig().LastBatch)
	for _, m := range forRollback {
		forRun := store.FindMigration(m.Name)
		if forRun == nil {
			log.Fatal("Migration", m.Name, "not found")
		}
		log.Println("Rolling back", forRun.GetName())
		forRun.Down(connection)
		log.Println("Rolled back", forRun.GetName())
	}
	model.RemoveLastBatch(connection, config.GetConfig().LastBatch)
}

func upOrIgnore(migration store.Migratable) {
	if !config.GetConfig().FirstRun && model.MigrationExist(connection, migration.GetName()) {
		return
	}
	log.Println("Migrating", migration.GetName())
	migration.Up(connection)
	model.AddMigrationRaw(connection, migration.GetName(), config.GetConfig().LastBatch+1)
	log.Println("Migrated", migration.GetName())
}

func parseFlags() {
	isRollback := flag.Bool("rollback", false, "Flag for init rollback.")
	envPath := flag.String("env-path", "", "Path to .env file.")
	envFile := flag.String("env-file", ".env", "Env file name.")
	verbose := flag.Bool("verbose", false, "Flag for show more info.")
	flag.Parse()
	config.GetConfig().IsRollback = *isRollback
	config.GetConfig().EnvPath = *envPath
	config.GetConfig().EnvFile = *envFile
	config.GetConfig().Verbose = *verbose
}

func prepare() {
	if config.GetConfig().Verbose {
		log.Println("load env file from:", config.GetConfig().GetEnvFullPath())
	}
	err := godotenv.Load(config.GetConfig().GetEnvFullPath())
	if err != nil {
		log.Println("Error loading .env file")
	}

	connector := db.NewConnector()
	connection, err = connector.Connect()
	if err != nil {
		log.Fatal(err)
	}
	config.GetConfig().LastBatch = model.GetLastBatch(connection)
	config.GetConfig().FirstRun = model.CreateMigrationsTable(connection)
}
