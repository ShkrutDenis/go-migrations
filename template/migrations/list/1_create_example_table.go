package list

import (
	"github.com/jmoiron/sqlx"
	"go-migrations/query_builder"
	"go-migrations/query_builders/mysql_builder"
)

type CreateExampleTable struct {}

func (m *CreateExampleTable) GetName() string {
	// Name of migration.
	// It will be stored to DB, must be unique.
	return "CreateExampleTable"
}

func (m *CreateExampleTable) Up(con *sqlx.DB) {
	// Write your migration logic here
	// Example:
	// 	con.MustExec("CREATE TABLE example ( id int auto_increment, constraint migrations_pk primary key (id));")
	table := mysql_builder.NewTable("test")
	table.String("field1", 255)

	con.MustExec(table.GetQuery())
}

func (m *CreateExampleTable) Down(con *sqlx.DB) {
	// Write your migration rollback logic here
	// Example:
	// 	con.MustExec("DROP TABLE example;")
}
