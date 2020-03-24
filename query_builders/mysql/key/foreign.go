package key

import (
	"fmt"
	"github.com/ShkrutDenis/go-migrations/query_builders/mysql/info"
)

type ForeignKey struct {
	name         string
	baseTable    string
	baseColumn   string
	targetTable  string
	targetColumn string
	onDelete     string
	onUpdate     string

	drop   bool
	change string
}

func NewForeignKey(table, baseColumn string) *ForeignKey {
	return &ForeignKey{baseTable: table, baseColumn: baseColumn, onDelete: "restrict", onUpdate: "restrict"}
}

func NewForeignKeyByKeyInfo(ki *info.KeyInfo) *ForeignKey {
	return &ForeignKey{
		name:         ki.ConstraintName,
		baseTable:    ki.TableName,
		baseColumn:   ki.ColumnName,
		targetTable:  ki.ReferencedTableName,
		targetColumn: ki.ReferencedColumnName,
	}
}

func (fk *ForeignKey) Reference(table string) *ForeignKey {
	fk.targetTable = table
	return fk
}

func (fk *ForeignKey) On(field string) *ForeignKey {
	fk.targetColumn = field
	return fk
}

func (fk *ForeignKey) OnUpdate(action string) *ForeignKey {
	fk.onUpdate = action
	return fk
}

func (fk *ForeignKey) OnDelete(action string) *ForeignKey {
	fk.onDelete = action
	return fk
}

func (fk *ForeignKey) Drop() *ForeignKey {
	fk.drop = true
	return fk
}

func (fk *ForeignKey) SetKeyName(name string) *ForeignKey {
	fk.name = name
	return fk
}

func (fk *ForeignKey) GenerateKeyName() *ForeignKey {
	fk.name = fmt.Sprintf("%v_%v_%v_fk", fk.baseTable, fk.targetTable, fk.targetColumn)
	return fk
}

func (fk *ForeignKey) GetSQL() string {
	if fk.drop {
		return fmt.Sprintf("ALTER TABLE %v drop foreign key %v;",
			fk.baseTable, fk.name) +
				fmt.Sprintf("Drop index %v on %v;",
			fk.name, fk.baseTable)
	}
	return fmt.Sprintf("constraint %v foreign key (%v) references %v (%v) on update %v on delete %v",
		fk.name, fk.baseColumn, fk.targetTable, fk.targetColumn, fk.onUpdate, fk.onDelete)
}

// Helpful functions
func (fk *ForeignKey) GetName() string {
	return fk.name
}

func (fk *ForeignKey) ForDrop() bool {
	return fk.drop
}
