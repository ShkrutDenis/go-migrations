package list

import (
	"github.com/jmoiron/sqlx"
	mysqlTable "go-migrations-local/query_builders/mysql_builder/table"
)

type UpdateExample1Table struct{}

func (m *UpdateExample1Table) GetName() string {
	// Name of migration.
	// It will be stored to DB, must be unique.
	return "UpdateExample1Table"
}

func (m *UpdateExample1Table) Up(con *sqlx.DB) {
	// Write your migration logic here
	// Example:
	// 	con.MustExec("CREATE TABLE example ( id int auto_increment, constraint migrations_pk primary key (id));")
	table := mysqlTable.ChangeTable("tournaments", con)
	table.Integer("test_fk")
	table.ForeignKey("test_fk").
		Reference("users").
		On("id").
		OnDelete("cascade").
		OnUpdate("cascade")
	table.DropForeignKey("tournaments_users_id_fk")
	table.MustExec()
}

func (m *UpdateExample1Table) Down(con *sqlx.DB) {
	// Write your migration rollback logic here
	// Example:
	// 	con.MustExec("DROP TABLE example;")
	table := mysqlTable.ChangeTable("tournaments", con)
	table.DropColumn("test_fk")
	table.ForeignKey("owner").
		Reference("users").
		On("id").
		OnDelete("cascade").
		OnUpdate("cascade")
	table.DropForeignKey("tournaments_users_id_fk")
	table.MustExec()
}
