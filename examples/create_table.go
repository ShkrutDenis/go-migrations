package list

import (
	"github.com/ShkrutDenis/go-migrations/builder"
	"github.com/jmoiron/sqlx"
)

type CreateExampleTable struct {}

func (m *CreateExampleTable) GetName() string {
	// Name of migration.
	// It will be stored to DB, must be unique.
	return "CreateExampleTable"
}

func (m *CreateExampleTable) Up(con *sqlx.DB) {
	// Create new table
	table := builder.NewTable("example", con)
	// Add primary key. It will be created column with type int and autoincrement.
	table.PrimaryKey("id")
	// Nullable column with default value. String(,n) => varchar(n).
	table.String("name", 100).Nullable().Default("value")
	// Builder has few predefined methods with type declaration.
	// Any way, you can use Column().Type() methods for create column with any type.
	table.Integer("count").Default("0")
	// Unique column
	table.Column("email").Type("varchar(50)").Unique()
	// Execute queries
	table.MustExec()
}

func (m *CreateExampleTable) Down(con *sqlx.DB) {
	// Drop table
	builder.DropTable("example", con).MustExec()
}
