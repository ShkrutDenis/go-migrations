package list

import (
	"github.com/ShkrutDenis/go-migrations/builder"
	"github.com/jmoiron/sqlx"
)

type ChangeExampleTable struct {}

func (m *ChangeExampleTable) GetName() string {
	// Name of migration.
	// It will be stored to DB, must be unique.
	return "ChangeExampleTable"
}

func (m *ChangeExampleTable) Up(con *sqlx.DB) {
	// Change existed table
	table := builder.ChangeTable("example", con)
	// For change column, just declare new column and call method Change().
	table.String("name", 100).Default("new value").NotNull().Change()
	// Remove unique condition
	table.Column("email").Type("varchar(50)").NotUnique().Change()
	// Add default timestamps (created_at nad updated_at)
	table.WithTimestamps()
	// Execute queries
	table.MustExec()
}

func (m *ChangeExampleTable) Down(con *sqlx.DB) {
	table := builder.ChangeTable("example", con)
	table.String("name", 100).Nullable().Default("value")
	table.Column("email").Type("varchar(50)").Unique()
	table.MustExec()
}
