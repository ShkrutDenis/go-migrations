package list

import (
	"github.com/jmoiron/sqlx"
	mysqlTable "go-migrations-local/query_builders/mysql_builder/table"
)

type CreateExample2Table struct{}

func (m *CreateExample2Table) GetName() string {
	// Name of migration.
	// It will be stored to DB, must be unique.
	return "CreateExample2Table"
}

func (m *CreateExample2Table) Up(con *sqlx.DB) {
	// Write your migration logic here
	// Example:
	// 	con.MustExec("CREATE TABLE example ( id int auto_increment, constraint migrations_pk primary key (id));")
	table := mysqlTable.NewTable("tournaments", con)
	table.Integer("id").Autoincrement()
	table.PrimaryKey("id")
	table.String("name", 255)
	table.Integer("owner")
	table.WithTimestamps()
	table.ForeignKey("owner").
		Reference("users").
		On("id").
		OnDelete("cascade").
		OnUpdate("cascade")
	table.MustExec()
}

func (m *CreateExample2Table) Down(con *sqlx.DB) {
	// Write your migration rollback logic here
	// Example:
	// 	con.MustExec("DROP TABLE example;")
	table := mysqlTable.DropTable("tournaments", con)
	table.MustExec()
}
