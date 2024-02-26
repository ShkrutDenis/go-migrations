package column

import (
	"fmt"
	"strings"

	"github.com/ShkrutDenis/go-migrations/builder/contract"
	"github.com/ShkrutDenis/go-migrations/builder/postgress/info"
	"github.com/jmoiron/sqlx"
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

	isPrimaryKey bool

	unique       bool
	hasUniqueKey bool

	drop   bool
	change bool
	first  bool
	after  string
	rename string

	info *info.ColumnInfo
}

func NewColumn(table, fieldName string, con *sqlx.DB) contract.Column {
	ci := info.GetColumnInfo(table, fieldName, con)
	c := &Column{name: fieldName, info: ci}
	c.init()
	return c
}

func (c *Column) init() contract.Column {
	if c.info != nil {
		c.fieldType = c.info.GetType()
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
func (c *Column) Type(fieldType string) contract.Column {
	c.fieldType = fieldType
	return c
}

func (c *Column) Nullable() contract.Column {
	c.nullable = null
	return c
}

func (c *Column) NotNull() contract.Column {
	c.nullable = notNull
	return c
}

func (c *Column) Autoincrement() contract.Column {
	c.autoincrement = true
	return c
}

func (c *Column) NotAutoincrement() contract.Column {
	c.autoincrement = false
	return c
}

func (c *Column) Default(value string) contract.Column {
	c.hasDefault = value != ""
	c.defaultValue = value
	return c
}

func (c *Column) Primary() contract.Column {
	c.isPrimaryKey = true
	return c
}

func (c *Column) Unique() contract.Column {
	c.unique = true
	return c
}

func (c *Column) NotUnique() contract.Column {
	c.unique = false
	return c
}

func (c *Column) Drop() contract.Column {
	c.drop = true
	return c
}

func (c *Column) Change() contract.Column {
	c.change = true
	return c
}

func (c *Column) First() contract.Column {
	c.first = true
	return c
}

func (c *Column) After(name string) contract.Column {
	c.after = name
	return c
}

func (c *Column) Rename(name string) contract.Column {
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
	return fmt.Sprintf("\"%s\"%s%s", c.name, c.columnOptionsSQL(), c.columnPositionSQL())
}

func (c *Column) changeColumnSQL() string {
	if c.info == nil {
		panic(fmt.Sprintf("Column %s not found. Column name is correct?", c.name))
	}
	var sql string
	if c.hasDefault {
		sql += fmt.Sprintf("ALTER COLUMN \"%s\" DROP DEFAULT,", c.name)
	}
	if c.fieldType != c.info.GetType() {
		sql += fmt.Sprintf("ALTER COLUMN \"%s\" TYPE %v USING \"%s\"::%v,", c.name, c.fieldType, c.name, c.fieldType)
	}
	if c.hasDefault {
		sql += fmt.Sprintf("ALTER COLUMN \"%s\" SET DEFAULT '%v',", c.name, c.defaultValue)
	}
	if c.info.Nullable() && c.nullable == notNull {
		sql += fmt.Sprintf("ALTER COLUMN \"%s\" SET NOT NULL,", c.name)
	} else if !c.info.Nullable() && c.nullable == null {
		sql += fmt.Sprintf("ALTER COLUMN \"%s\" DROP NOT NULL,", c.name)
	}
	return strings.TrimRight(sql, ",")
}

func (c *Column) renameColumnSQL() string {
	return fmt.Sprintf("rename column \"%s\" to \"%s\";", c.name, c.rename)
}

func (c *Column) dropColumnSQL() string {
	return fmt.Sprintf("drop column \"%s\",", c.name)
}

func (c *Column) columnOptionsSQL() string {
	sql := " " + c.fieldType

	if c.nullable == "" {
		c.nullable = notNull
	}
	sql += " " + c.nullable

	if c.hasDefault {
		if string(c.defaultValue[len(c.defaultValue)-1:]) == ")" {
			sql += fmt.Sprintf(" DEFAULT %v", c.defaultValue)
		} else {
			sql += fmt.Sprintf(" DEFAULT '%v'", c.defaultValue)
		}
	}

	if c.isPrimaryKey {
		if strings.ToLower(c.fieldType) == "bigint" {
			return " BIGSERIAL PRIMARY KEY"
		} else if strings.Contains(c.fieldType, "int") {
			return " SERIAL PRIMARY KEY"
		}
		sql += " PRIMARY KEY"
	}

	return sql
}

func (c *Column) columnPositionSQL() string {
	// TODO:Postgres don`t support positioning
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

func (c *Column) IsPrimary() bool {
	return c.isPrimaryKey
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
