package key

import "fmt"

type PrimaryKey struct {
	name  string
	field string
	table string

	drop bool
}

func NewPrimaryKey(table, field string) *PrimaryKey {
	return &PrimaryKey{field: field, table: table}
}

func (pk *PrimaryKey) Drop() *PrimaryKey {
	pk.drop = true
	return pk
}

func (pk *PrimaryKey) SetKeyName(name string) *PrimaryKey {
	pk.name = name
	return pk
}

func (pk *PrimaryKey) GenerateKeyName() *PrimaryKey {
	pk.name = fmt.Sprintf("%v_pk", pk.table)
	return pk
}

func (pk *PrimaryKey) GetSQL() string {
	if pk.drop {
		return fmt.Sprintf("Drop index %v on %v;",
			pk.name, pk.table)
	}
	return fmt.Sprintf("constraint %v primary key (%v)", pk.name, pk.field)
}

// Helpful functions
func (pk *PrimaryKey) GetName() string {
	return pk.name
}
