package list

import (
	"github.com/jmoiron/sqlx"
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
	// Or you can use existed query builder
}

func (m *CreateExampleTable) Down(con *sqlx.DB) {
	// Write your migration rollback logic here
	// Example:
	// 	con.MustExec("DROP TABLE example;")
	// Or you can use existed query builder
}
