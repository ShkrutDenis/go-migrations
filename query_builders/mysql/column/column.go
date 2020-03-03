package column

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go-migrations-local/query_builders/mysql/info"
)

const (
	null    = "NULL"
	notNull = "NOT NULL"
)

type Column struct {
	name      string
	fieldType string

	hasDefault   bool
	defaultValue string

	nullable      string
	autoincrement bool

	unique       bool
	hasUniqueKey bool

	drop   bool
	change bool
	first  bool
	after  string
	rename string

	info *info.ColumnInfo
}

func NewColumn(table, fieldName string, con *sqlx.DB) *Column {
	ci := info.GetColumnInfo(table, fieldName, con)
	c := &Column{name: fieldName, info: ci}
	c.Init()
	return c
}

func (c *Column) Init() *Column {
	if c.info != nil {
		c.fieldType = c.info.ColumnType
		if c.info.Nullable() {
			c.nullable = null
		} else {
			c.nullable = notNull
		}
		c.unique = c.info.IsUnique()
		c.hasUniqueKey = c.unique
		c.hasDefault, c.defaultValue = c.info.HasDefault()
		c.autoincrement = c.info.IsPrimary()
	}
	return c
}

// Functions for modify table
func (c *Column) Type(fieldType string) *Column {
	c.fieldType = fieldType
	return c
}

func (c *Column) Nullable() *Column {
	c.nullable = null
	return c
}

func (c *Column) NotNull() *Column {
	c.nullable = notNull
	return c
}

func (c *Column) Autoincrement() *Column {
	c.autoincrement = true
	return c
}

func (c *Column) NotAutoincrement() *Column {
	c.autoincrement = false
	return c
}

func (c *Column) Default(value string) *Column {
	c.hasDefault = value != ""
	c.defaultValue = value
	return c
}

func (c *Column) Unique() *Column {
	c.unique = true
	return c
}

func (c *Column) NotUnique() *Column {
	c.unique = false
	return c
}

func (c *Column) Drop() *Column {
	c.drop = true
	return c
}

func (c *Column) Change() *Column {
	c.change = true
	return c
}

func (c *Column) First() *Column {
	c.first = true
	return c
}

func (c *Column) After(name string) *Column {
	c.after = name
	return c
}

func (c *Column) Rename(name string) *Column {
	c.rename = name
	return c
}

// Functions for generate SQL
func (c *Column) GetSQL() string {
	if c.drop {
		return c.dropColumnSQL()
	}
	if c.rename != "" {
		return c.renameColumnSQL()
	}
	if c.change {
		return c.changeColumnSQL()
	}
	return c.addColumnSQL()
}

func (c *Column) addColumnSQL() string {
	return c.name + c.columnOptionsSQL() + c.columnPositionSQL()
}

func (c *Column) changeColumnSQL() string {
	sql := fmt.Sprintf("modify column %v", c.name)
	sql += c.columnOptionsSQL()
	sql += c.columnPositionSQL()
	return sql
}

func (c *Column) renameColumnSQL() string {
	sql := fmt.Sprintf("change column %v %v", c.name, c.rename)
	sql += c.columnOptionsSQL()
	sql += c.columnPositionSQL()
	return sql + ";"
}

func (c *Column) dropColumnSQL() string {
	return fmt.Sprintf("drop column %v;", c.name)
}

func (c *Column) columnOptionsSQL() string {
	sql := " " + c.fieldType
	if c.hasDefault {
		sql += " default " + c.defaultValue
	}
	if c.autoincrement {
		sql += " auto_increment"
	} else {
		if c.nullable == "" {
			c.nullable = notNull
		}
		sql += " " + c.nullable
	}
	return sql
}

func (c *Column) columnPositionSQL() string {
	if c.first {
		return " first"
	}
	if c.after != "" {
		return " after " + c.after
	}
	return ""
}

// Helpful functions
func (c *Column) GetName() string {
	return c.name
}

func (c *Column) GetUniqueKeyName() string {
	if c.info == nil {
		return ""
	}
	k := c.info.GetUniqueKey()
	if k == nil {
		return ""
	}
	return k.KeyName
}

func (c *Column) IsUnique() bool {
	return c.unique
}

func (c *Column) HasUniqueKey() bool {
	return c.hasUniqueKey
}

func (c *Column) NeedUniqueKey() bool {
	return c.unique && !c.hasUniqueKey
}

func (c *Column) NeedDropUniqueKey() bool {
	return !c.unique && c.hasUniqueKey
}

func (c *Column) IsWaitingDrop() bool {
	return c.drop
}

func (c *Column) IsWaitingRename() bool {
	return c.rename != ""
}

func (c *Column) IsWaitingChange() bool {
	return c.change
}
