package engine

import (
	"fmt"
	"github.com/ShkrutDenis/go-migrations/engine/config"
	"github.com/ShkrutDenis/go-migrations/engine/db"
	"github.com/ShkrutDenis/go-migrations/engine/migrator"
	"github.com/ShkrutDenis/go-migrations/engine/store"
	"log"
)

type Engine struct {
	migrator   *migrator.Migrator
	migrations store.Migrations

	config    *config.Config
	connector *db.Connector
}

func NewEngine() *Engine {
	return &Engine{
		migrations: store.Migrations{},
		config:     config.NewConfig(),
		connector:  db.NewConnector(),
	}
}

//GetConnector returns current DB connector
func (e *Engine) GetConnector() *db.Connector {
	return e.connector
}

//WithConfig overrides config for the migrator
func (e *Engine) WithConfig(c config.Config) {
	e.config = config.NewFromConfig(c)
}

func (e *Engine) Run(migs []store.Migratable) {
	e.init()
	e.migrations.Register(migs)

	if e.config.IsRollback {
		e.rollback()
		return
	}

	for _, m := range migs {
		e.upOrIgnore(m)
	}
	e.CheckExtraMigrations()
}

func (e *Engine) rollback() {
	forRollback := e.migrator.GetLastMigrations()
	for _, m := range forRollback {
		forRun := e.migrations.Find(m.Name)
		if forRun == nil {
			log.Fatalf("Migration %s not found", m.Name)
		}
		log.Println("Rolling back", forRun.GetName())
		forRun.Down(e.connector.GetConnection())
		log.Println("Rolled back", forRun.GetName())
	}
	e.migrator.RemoveLastBatch()
}

func (e *Engine) upOrIgnore(migration store.Migratable) {
	if !e.config.FirstRun && e.migrator.MigrationExist(migration.GetName()) {
		return
	}
	log.Println("Migrating", migration.GetName())
	migration.Up(e.connector.GetConnection())
	e.migrator.AddMigrationRaw(migration.GetName())
	log.Println("Migrated", migration.GetName())
}

//CheckExtraMigrations checks mismatching between list of migrations and migrations from DB
func (e *Engine) CheckExtraMigrations() {
	extra := e.migrator.GetExtraMigrations(e.migrations.GetNames())
	if len(extra) != 0 {
		logMsg := "Found extra migrations in DB.\nMigrations:"
		for _, m := range extra {
			logMsg = fmt.Sprintf("%s\n - %s", logMsg, m.Name)
		}
		log.Println(logMsg)
	}
}

func (e *Engine) init() {
	e.config.InitByFlags()
	e.config.InitEnv()

	if err := e.connector.Connect(); err != nil {
		log.Fatalln("Can't create DB connection. Error: ", err)
	}
	e.migrator = migrator.NewMigrator(e.connector.GetConnection(), e.connector.GetDriver())

	e.config.FirstRun = e.migrator.CreateMigrationsTable()
}
