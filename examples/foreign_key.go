package list

import (
	"github.com/ShkrutDenis/go-migrations/builder"
	"github.com/jmoiron/sqlx"
)

type ForeignKeyExampleTable struct {}

func (m *ForeignKeyExampleTable) GetName() string {
	// Name of migration.
	// It will be stored to DB, must be unique.
	return "ForeignKeyExampleTable"
}

func (m *ForeignKeyExampleTable) Up(con *sqlx.DB) {
	// Create new table
	table := builder.NewTable("fk_example", con)
	// Add primary key. It will be created column with type int and autoincrement.
	table.PrimaryKey("id")
	// Add foreign key. First you need add column then add foreign key.
	table.Integer("example_id")
	table.ForeignKey("example_id").Reference("example").On("id").OnDelete("cascade").OnUpdate("cascade")
	// Execute queries
	table.MustExec()
}

func (m *ForeignKeyExampleTable) Down(con *sqlx.DB) {
	// Drop table
	builder.DropTable("fk_example", con).MustExec()
}
