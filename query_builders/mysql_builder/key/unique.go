package key

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UniqueKey struct {
	name  string
	table string
	field string

	drop bool
}

func NewUniqueKey(table, field string) *UniqueKey {
	uk := &UniqueKey{table: table, field: field}
	return uk
}

func (uk *UniqueKey) SetKeyName(name string) *UniqueKey {
	uk.name = name
	return uk
}

func (uk *UniqueKey) GenerateKeyName() *UniqueKey {
	uk.name = fmt.Sprintf("%v_%v_uindex", uk.table, uk.field)
	return uk
}

func (uk *UniqueKey) Drop() *UniqueKey {
	uk.drop = true
	return uk
}

func (uk *UniqueKey) GetSQL() string {
	if uk.drop {
		return fmt.Sprintf("Drop index %v on %v;",
			uk.name, uk.table)
	}
	return fmt.Sprintf("Create unique index %v on %v (%v);",
		uk.name, uk.table, uk.field)
}

func (uk *UniqueKey) Exec(con *sqlx.DB) error {
	_, err := con.Exec(uk.GetSQL())
	return err
}

func (uk *UniqueKey) MustExec(con *sqlx.DB) {
	con.MustExec(uk.GetSQL())
}

// Helpful functions
func (uk *UniqueKey) GetName() string {
	return uk.name
}
