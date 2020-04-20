package key

import (
	"fmt"
	"github.com/ShkrutDenis/go-migrations/builder/contract"
	"github.com/ShkrutDenis/go-migrations/builder/postgress/info"
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

func NewForeignKey(table, baseColumn string) contract.ForeignKey {
	return &ForeignKey{baseTable: table, baseColumn: baseColumn, onDelete: "NO ACTION", onUpdate: "NO ACTION"}
}

func NewForeignKeyByKeyInfo(ki *info.KeyInfo) contract.ForeignKey {
	return &ForeignKey{
		name:         ki.ConstraintName,
		baseTable:    ki.BaseTable,
		baseColumn:   ki.BaseColumn,
		targetTable:  ki.TargetTable,
		targetColumn: ki.TargetColumn,
		onUpdate:     ki.OnUpdate,
		onDelete:     ki.OnDelete,
	}
}

func (fk *ForeignKey) Reference(table string) contract.ForeignKey {
	fk.targetTable = table
	return fk
}

func (fk *ForeignKey) On(field string) contract.ForeignKey {
	fk.targetColumn = field
	return fk
}

func (fk *ForeignKey) OnUpdate(action string) contract.ForeignKey {
	fk.onUpdate = action
	return fk
}

func (fk *ForeignKey) OnDelete(action string) contract.ForeignKey {
	fk.onDelete = action
	return fk
}

func (fk *ForeignKey) Drop() contract.ForeignKey {
	fk.drop = true
	return fk
}

func (fk *ForeignKey) SetKeyName(name string) contract.ForeignKey {
	fk.name = name
	return fk
}

func (fk *ForeignKey) GenerateKeyName() contract.ForeignKey {
	fk.name = fmt.Sprintf("%v_%v_%v_fkey", fk.baseTable, fk.targetTable, fk.targetColumn)
	return fk
}

func (fk *ForeignKey) GetSQL() string {
	if fk.drop {
		return fmt.Sprintf("ALTER TABLE %v DROP CONSTRAINT %v;",
			fk.baseTable, fk.name)
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
