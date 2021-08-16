package key

import (
	"fmt"
	"strings"
)

type PrimaryKey struct {
	name  string
	table string

	columns []string

	autoincrement map[string]bool
	drop          map[string]bool
}

func NewPrimaryKey(table string) *PrimaryKey {
	return &PrimaryKey{table: table, autoincrement: make(map[string]bool), drop: make(map[string]bool)}
}

func (pk *PrimaryKey) Drop(column string) *PrimaryKey {
	pk.drop[column] = true
	return pk
}

func (pk *PrimaryKey) AddColumn(column string, autoincrement bool) *PrimaryKey {
	pk.columns = append(pk.columns, column)
	pk.autoincrement[column] = autoincrement
	return pk
}

func (pk *PrimaryKey) SetKeyName(name string) *PrimaryKey {
	pk.name = name
	return pk
}

func (pk *PrimaryKey) GenerateKeyName() *PrimaryKey {
	pk.name = fmt.Sprintf("%v_%v_pk", pk.table, strings.Join(pk.columns, "_"))
	return pk
}

func (pk *PrimaryKey) GenerateSequenceName(column string) string {
	return fmt.Sprintf("%v_%v_seq", pk.table, column)
}

func (pk *PrimaryKey) GetSQL() string {
	if len(pk.columns) == 0 {
		return ""
	}
	if len(pk.columns) > 1 {
		return fmt.Sprintf("ALTER TABLE %v ADD PRIMARY KEY (%v);", pk.table, strings.Join(pk.columns, ","))
	}
	column := pk.columns[0]
	if pk.drop[column] {
		return fmt.Sprintf("DROP INDEX %v ON %v;", pk.GetName(), pk.table)
	}
	sql := fmt.Sprintf("ALTER TABLE %v ", pk.table)
	if pk.autoincrement[column] {
		seq := pk.GenerateSequenceName(column)
		sql = fmt.Sprintf("CREATE SEQUENCE %v OWNED BY %v.%v;", seq, pk.table, column) +
			sql + fmt.Sprintf("ALTER COLUMN %v SET DEFAULT nextval('%v'),", column, seq)
	}

	return sql + fmt.Sprintf("ADD PRIMARY KEY (%v);", column)
}

// Helpful functions
func (pk *PrimaryKey) GetName() string {
	if pk.name == "" {
		pk.GenerateKeyName()
	}
	return pk.name
}
