package list

import (
	"github.com/jmoiron/sqlx"
	mysqlTable "go-migrations-local/query_builders/mysql_builder/table"
)

type CreateExample1Table struct{}

func (m *CreateExample1Table) GetName() string {
	// Name of migration.
	// It will be stored to DB, must be unique.
	return "CreateExample1Table"
}

func (m *CreateExample1Table) Up(con *sqlx.DB) {
	// Write your migration logic here
	// Example:
	// 	con.MustExec("CREATE TABLE example ( id int auto_increment, constraint migrations_pk primary key (id));")
	table := mysqlTable.NewTable("users", con)
	table.Integer("id").Autoincrement()
	table.PrimaryKey("id")
	table.String("first_name", 255)
	table.String("last_name", 255)
	table.String("email", 255).Unique()
	table.WithTimestamps()
	table.MustExec()
}

func (m *CreateExample1Table) Down(con *sqlx.DB) {
	// Write your migration rollback logic here
	// Example:
	// 	con.MustExec("DROP TABLE example;")
	table := mysqlTable.DropTable("users", con)
	table.MustExec()
}
