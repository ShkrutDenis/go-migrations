package key

import (
	"fmt"
	"github.com/ShkrutDenis/go-migrations/builder/contract"
	"github.com/ShkrutDenis/go-migrations/engine/config"
	"github.com/jmoiron/sqlx"
	"log"
)

type UniqueKey struct {
	name  string
	table string
	field string

	drop bool
}

func NewUniqueKey(table, field string) contract.UniqueKey {
	uk := &UniqueKey{table: table, field: field}
	return uk
}

func (uk *UniqueKey) SetKeyName(name string) contract.UniqueKey {
	uk.name = name
	return uk
}

func (uk *UniqueKey) GenerateKeyName() contract.UniqueKey {
	uk.name = fmt.Sprintf("unique_%v_%v", uk.table, uk.field)
	return uk
}

func (uk *UniqueKey) Drop() contract.UniqueKey {
	uk.drop = true
	return uk
}

func (uk *UniqueKey) GetSQL() string {
	if uk.drop {
		return fmt.Sprintf("ALTER TABLE %v DROP CONSTRAINT %v;",
			uk.table, uk.name)
	}
	return fmt.Sprintf("Create unique index %v_%v on %v (%v);",
		uk.table, uk.field, uk.table, uk.field) +
		fmt.Sprintf(" ALter table %v add constraint %v unique using index %v_%v;",
			uk.table, uk.name, uk.table, uk.field)
}

func (uk *UniqueKey) Exec(con *sqlx.DB) error {
	if config.Verbose {
		log.Println(uk.GetSQL())
	}
	_, err := con.Exec(uk.GetSQL())
	return err
}

func (uk *UniqueKey) MustExec(con *sqlx.DB) {
	if config.Verbose {
		log.Println(uk.GetSQL())
	}
	con.MustExec(uk.GetSQL())
}

// Helpful functions
func (uk *UniqueKey) GetName() string {
	return uk.name
}
