package migrator

import (
	"github.com/jmoiron/sqlx"
)

type Migrator struct {
	driver string
	db     *sqlx.DB

	lastBatch int
}

func NewMigrator(db *sqlx.DB, driver string) *Migrator {
	m := &Migrator{db: db, driver: driver}
	m.lastBatch = m.GetLastBatch()
	return m
}

func (m *Migrator) CreateMigrationsTable() bool {
	_, err := m.db.Exec(m.checkTableExitSQL())
	if err != nil {
		m.db.MustExec(m.creteTableSQL())
		return true
	}
	return false
}

func (m *Migrator) GetExtraMigrations(names []string) []Migration {
	var list []Migration
	q, a, e := sqlx.In(m.checkExtraRawsSQL(), names)
	if e != nil {
		panic(e)
	}
	q = m.db.Rebind(q)
	e = m.db.Select(&list, q, a...)
	if e != nil {
		panic(e)
	}
	return list
}

func (m *Migrator) AddMigrationRaw(migration string) {
	m.db.MustExec(m.addRawSQL(), migration, m.lastBatch+1)
}

func (m *Migrator) MigrationExist(name string) bool {
	var raw Migration
	err := m.db.Get(&raw, m.getRawSQL(), name)
	if err != nil {
		return false
	}
	return true
}

func (m *Migrator) GetLastMigration() (bool, *Migration) {
	var raw Migration
	err := m.db.Get(&raw, m.getLastRawSQL())
	if err != nil {
		return false, &raw
	}
	return true, &raw
}

func (m *Migrator) GetLastMigrations() []*Migration {
	var list []*Migration
	err := m.db.Select(&list, m.getLastRawsSQL(), m.lastBatch)
	if err != nil {
		panic(err)
	}
	return list
}

func (m *Migrator) RemoveMigrationRaw(migration string) {
	m.db.MustExec(m.removeMigrationSQL(), migration)
}

func (m *Migrator) GetLastBatch() int {
	var raw Migration
	_ = m.db.Get(&raw, m.getLastBatchSQL())
	return raw.Batch
}

func (m *Migrator) RemoveLastBatch() {
	m.db.MustExec(m.removeBatchSQL(), m.lastBatch)
}

// SQL queries functions
func (m *Migrator) checkTableExitSQL() string {
	switch m.driver {
	case "mysql":
		return "SELECT * FROM migrations LIMIT 1;"
	case "postgres":
		return "SELECT * FROM migrations LIMIT 1;"
	default:
		panic("Not supported DB driver: " + m.driver)
	}
}

func (m *Migrator) checkExtraRawsSQL() string {
	switch m.driver {
	case "mysql":
		return "SELECT * FROM migrations WHERE name NOT IN (?);"
	case "postgres":
		return "SELECT * FROM migrations WHERE name NOT IN (?);"
	default:
		panic("Not supported DB driver: " + m.driver)
	}
}

func (m *Migrator) creteTableSQL() string {
	switch m.driver {
	case "mysql":
		return "CREATE TABLE migrations (id int auto_increment, name varchar(255) not null, batch int not null, created_at datetime default current_timestamp not null, constraint migrations_pk primary key (id));"
	case "postgres":
		return "CREATE TABLE migrations (id serial PRIMARY KEY, name varchar(255) not null, batch int not null, created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP not null);"
	default:
		panic("Not supported DB driver: " + m.driver)
	}
}

func (m *Migrator) getRawSQL() string {
	switch m.driver {
	case "mysql":
		return "SELECT * FROM migrations WHERE name=?;"
	case "postgres":
		return "SELECT * FROM migrations WHERE name=$1;"
	default:
		panic("Not supported DB driver: " + m.driver)
	}
}

func (m *Migrator) addRawSQL() string {
	switch m.driver {
	case "mysql":
		return "INSERT INTO migrations (name, batch) VALUES (?, ?);"
	case "postgres":
		return "INSERT INTO migrations (name, batch) VALUES ($1, $2);"
	default:
		panic("Not supported DB driver: " + m.driver)
	}
}

func (m *Migrator) getLastRawSQL() string {
	switch m.driver {
	case "mysql":
		return "SELECT * FROM migrations ORDER BY created_at DESC, id DESC LIMIT 1;"
	case "postgres":
		return "SELECT * FROM migrations ORDER BY created_at DESC, id DESC LIMIT 1;"
	default:
		panic("Not supported DB driver: " + m.driver)
	}
}

func (m *Migrator) getLastBatchSQL() string {
	switch m.driver {
	case "mysql":
		return "SELECT * FROM migrations ORDER BY batch DESC LIMIT 1;"
	case "postgres":
		return "SELECT * FROM migrations ORDER BY batch DESC LIMIT 1;"
	default:
		panic("Not supported DB driver: " + m.driver)
	}
}

func (m *Migrator) getLastRawsSQL() string {
	switch m.driver {
	case "mysql":
		return "SELECT * FROM migrations WHERE batch=? ORDER BY created_at DESC, id DESC;"
	case "postgres":
		return "SELECT * FROM migrations WHERE batch=$1 ORDER BY created_at DESC, id DESC;"
	default:
		panic("Not supported DB driver: " + m.driver)
	}
}

func (m *Migrator) removeBatchSQL() string {
	switch m.driver {
	case "mysql":
		return "DELETE FROM migrations WHERE batch=?;"
	case "postgres":
		return "DELETE FROM migrations WHERE batch=$1;"
	default:
		panic("Not supported DB driver: " + m.driver)
	}
}

func (m *Migrator) removeMigrationSQL() string {
	switch m.driver {
	case "mysql":
		return "DELETE FROM migrations WHERE name=$1"
	case "postgres":
		return "DELETE FROM migrations WHERE name=$1"
	default:
		panic("Not supported DB driver: " + m.driver)
	}
}
