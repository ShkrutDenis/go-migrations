package table

import (
	"fmt"
	"github.com/ShkrutDenis/go-migrations/query_builders/mysql_builder/column"
	"github.com/ShkrutDenis/go-migrations/query_builders/mysql_builder/info"
	"github.com/ShkrutDenis/go-migrations/query_builders/mysql_builder/key"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type Table struct {
	name string

	primaryKey  *key.PrimaryKey
	foreignKeys []*key.ForeignKey
	uniqueKeys  []*key.UniqueKey

	columns    []*column.Column
	timestamps bool

	drop   bool
	change bool

	connect *sqlx.DB
}

func NewTable(name string, con *sqlx.DB) *Table {
	return &Table{name: name, connect: con.Unsafe()}
}

func DropTable(name string, con *sqlx.DB) *Table {
	return &Table{name: name, drop: true, connect: con.Unsafe()}
}

func ChangeTable(name string, con *sqlx.DB) *Table {
	return &Table{name: name, change: true, connect: con.Unsafe()}
}

// Functions for table columns
func (t *Table) Column(name string) *column.Column {
	c := column.NewColumn(t.name, name, t.connect)
	t.columns = append(t.columns, c)
	return c
}

func (t *Table) String(name string, length int) *column.Column {
	if length < 0 {
		length = 255
	}
	c := t.Column(name).Type("varchar(" + strconv.Itoa(length) + ")")
	return c
}

func (t *Table) Integer(name string) *column.Column {
	c := t.Column(name).Type("int")
	return c
}

func (t *Table) WithTimestamps() *Table {
	t.timestamps = true
	return t
}

func (t *Table) getTimeStampsSQL() string {
	return "created_at datetime default current_timestamp not null," +
		"updated_at datetime default current_timestamp ON UPDATE current_timestamp not null"
}

func (t *Table) RenameColumn(oldName, newName string) *column.Column {
	c := t.Column(oldName).Rename(newName)
	return c
}

func (t *Table) DropColumn(name string) *column.Column {
	c := t.Column(name).Drop()
	return c
}

// Functions for keys
func (t *Table) PrimaryKey(Column string) *Table {
	t.primaryKey = key.NewPrimaryKey(t.name, Column).GenerateKeyName()
	return t
}

func (t *Table) ForeignKey(Column string) *key.ForeignKey {
	k := key.NewForeignKey(t.name, Column)
	t.foreignKeys = append(t.foreignKeys, k)
	return k
}

func (t *Table) DropForeignKey(name string) *Table {
	if !t.change {
		return t
	}
	ki := info.GetKeyInfo(name, t.connect)
	if ki == nil {
		panic("Foreign key " + name + " not exist")
	}
	t.foreignKeys = append(t.foreignKeys, key.NewForeignKeyByKeyInfo(ki).Drop())
	return t
}

// Generate SQL functions
func (t *Table) GetSQL() string {
	if t.drop {
		return t.dropTableSQL()
	}
	sql := ""
	for _, k := range t.foreignKeys {
		if k.ForDrop() {
			sql += k.GetSQL()
		}
	}
	if t.change {
		return sql + t.changeTableSQL()
	}
	return sql + t.createTableSQL()
}

func (t *Table) createTableSQL() string {
	t.uniqueKeys = t.uniqueKeys[:0]
	sql := "CREATE TABLE " + t.name + "("
	for _, c := range t.columns {
		sql += c.GetSQL() + ","
		t.checkUniqueKey(c)
	}
	if t.timestamps {
		sql += t.getTimeStampsSQL() + ","
	}
	for _, k := range t.foreignKeys {
		if k.ForDrop() {
			continue
		}
		if k.GetName() == "" {
			k.GenerateKeyName()
		}
		sql += k.GetSQL() + ","
	}
	sql += t.primaryKey.GetSQL()
	return sql + ");"
}

func (t *Table) changeTableSQL() string {
	base := "ALTER TABLE " + t.name + " "
	var forAdd string
	var forModify string
	var forRename string
	var forDrop string
	var fKeys string
	for _, c := range t.columns {
		if c.IsWaitingDrop() {
			forDrop += c.GetSQL()
			continue
		}
		t.checkUniqueKey(c)
		if c.IsWaitingRename() {
			forRename += c.GetSQL()
			continue
		}
		if c.IsWaitingChange() {
			forModify += c.GetSQL() + ","
			continue
		}
		forAdd += "add " + c.GetSQL() + ","
	}
	for _, k := range t.foreignKeys {
		if k.ForDrop() {
			continue
		}
		if k.GetName() == "" {
			k.GenerateKeyName()
		}
		fKeys += "add " + k.GetSQL() + ","
	}
	if forDrop != "" {
		forDrop = base + forDrop
	}
	if forRename != "" {
		forRename = base + forRename + ";"
	}
	if forModify != "" {
		forModify = base + strings.TrimRight(forModify, ",") + ";"
	}
	if forAdd != "" {
		forAdd = base + strings.TrimRight(forAdd, ",") + ";"
	}
	if fKeys != "" {
		fKeys = base + strings.TrimRight(fKeys, ",") + ";"
	}
	return forRename + forDrop + forModify + forAdd + fKeys
}

func (t *Table) dropTableSQL() string {
	return fmt.Sprintf("DROP TABLE %v;", t.name)
}

// Execution functions
func (t *Table) Exec() error {
	queries := strings.Split(t.GetSQL(), ";")
	for i, q := range queries {
		if i == len(queries) {
			break
		}
		_, err := t.connect.Exec(q + ";")
		if err != nil {
			return err
		}
	}

	for _, k := range t.uniqueKeys {
		if err := k.Exec(t.connect); err != nil {
			return err
		}
	}

	return nil
}

func (t *Table) MustExec() {
	queries := strings.Split(t.GetSQL(), ";")
	for _, q := range queries {
		if q == "" {
			break
		}
		t.connect.MustExec(q + ";")
	}

	for _, k := range t.uniqueKeys {
		k.MustExec(t.connect)
	}
}

// Helpful functions
func (t *Table) checkUniqueKey(c *column.Column) {
	if c.NeedUniqueKey() {
		k := key.NewUniqueKey(t.name, c.GetName()).GenerateKeyName()
		t.uniqueKeys = append(t.uniqueKeys, k)
	}
	if c.NeedDropUniqueKey() {
		k := key.NewUniqueKey(t.name, c.GetName()).SetKeyName(c.GetUniqueKeyName()).Drop()
		t.uniqueKeys = append(t.uniqueKeys, k)
	}
}
