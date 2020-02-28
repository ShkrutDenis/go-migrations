package mysql_builder

import "fmt"

type PrimaryKey struct {
	name  string
	field string
}

func NewPrimaryKey(name, field string) *PrimaryKey {
	return &PrimaryKey{name: name, field: field}
}

func (pk *PrimaryKey) GetSQL() string {
	return fmt.Sprintf("constraint %v primary key (%v)", pk.name, pk.field)
}

type ForeignKey struct {
	name        string
	baseField   string
	targetTable string
	targetField string
	onDelete    string
	onUpdate    string
}

func NewForeignKey(name, baseField string) *ForeignKey {
	return &ForeignKey{name: name, baseField: baseField}
}

func (fk *ForeignKey) Reference(table string) *ForeignKey {
	fk.targetTable = table
	return fk
}

func (fk *ForeignKey) On(field string) *ForeignKey {
	fk.targetField = field
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

func (fk *ForeignKey) GetSQL() string {
	return fmt.Sprintf("constraint %v foreign key (%v) references %v (%v) on update %v on delete %v",
		fk.name, fk.baseField, fk.targetTable, fk.targetField, fk.onUpdate, fk.onDelete)
}
