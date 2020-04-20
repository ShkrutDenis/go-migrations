package go_migrations

import (
	"flag"
	"github.com/ShkrutDenis/go-migrations/db"
	"github.com/ShkrutDenis/go-migrations/model"
	"github.com/ShkrutDenis/go-migrations/store"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var connection *sqlx.DB
var config *Config

type Config struct {
	IsRollback bool
	EnvPath    string
	EnvFile    string

	lastBatch int
	firstRun  bool
}

func (c *Config) GetEnvFullPath() string {
	if c.EnvPath == "" {
		return c.EnvFile
	}
	_, err := os.Stat(c.EnvPath + "/" + c.EnvFile)
	if os.IsNotExist(err) {
		return c.EnvPath + "\\" + c.EnvFile
	}
	return c.EnvPath + "/" + c.EnvFile
}

func init() {
	config = &Config{}
}

func Run(migs []store.Migratable) {
	parseFlags()
	prepare()

	if config.IsRollback {
		store.RegisterMigrations(migs)
		rollBack()
		return
	}

	for _, m := range migs {
		upOrIgnore(m)
	}
}

func rollBack() {
	forRollback := model.GetLastMigrations(connection, config.lastBatch)
	for _, m := range forRollback {
		forRun := store.FindMigration(m.Name)
		if forRun == nil {
			log.Fatal("Migration", m.Name, "not found")
		}
		log.Println("Rolling back", forRun.GetName())
		forRun.Down(connection)
		log.Println("Rolled back", forRun.GetName())
	}
	model.RemoveLastBatch(connection, config.lastBatch)
}

func upOrIgnore(migration store.Migratable) {
	if !config.firstRun && model.MigrationExist(connection, migration.GetName()) {
		return
	}
	log.Println("Migrating", migration.GetName())
	migration.Up(connection)
	model.AddMigrationRaw(connection, migration.GetName(), config.lastBatch+1)
	log.Println("Migrated", migration.GetName())
}

func parseFlags() {
	isRollback := flag.Bool("rollback", false, "Flag for init rollback.")
	envPath := flag.String("env-path", "", "Path to .env file.")
	envFile := flag.String("env-file", ".env", "Env file name.")
	flag.Parse()
	config.IsRollback = *isRollback
	config.EnvPath = *envPath
	config.EnvFile = *envFile
}

func prepare() {
	err := godotenv.Load(config.GetEnvFullPath())
	if err != nil {
		log.Println("Error loading .env file")
	}

	connector := db.NewConnector()
	connection, err = connector.Connect()
	if err != nil {
		log.Fatal(err)
	}
	config.lastBatch = model.GetLastBatch(connection)
	config.firstRun = model.CreateMigrationsTable(connection)
}
