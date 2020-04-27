package table

import (
	"fmt"
	"github.com/ShkrutDenis/go-migrations/builder/contract"
	"github.com/ShkrutDenis/go-migrations/builder/postgress/column"
	"github.com/ShkrutDenis/go-migrations/builder/postgress/info"
	"github.com/ShkrutDenis/go-migrations/builder/postgress/key"
	"github.com/ShkrutDenis/go-migrations/config"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
	"strings"
)

type Table struct {
	name    string
	newName string

	foreignKeys []contract.ForeignKey
	uniqueKeys  []contract.UniqueKey

	columns    []contract.Column
	timestamps bool

	drop   bool
	change bool

	connect *sqlx.DB
}

func NewTable(name string, con *sqlx.DB) contract.Table {
	return &Table{name: name, connect: con.Unsafe()}
}

func DropTable(name string, con *sqlx.DB) contract.Table {
	return &Table{name: name, drop: true, connect: con.Unsafe()}
}

func ChangeTable(name string, con *sqlx.DB) contract.Table {
	return &Table{name: name, change: true, connect: con.Unsafe()}
}

func RenameTable(oldName, newName string, con *sqlx.DB) contract.Table {
	return &Table{name: oldName, newName: newName, connect: con.Unsafe()}
}

// Functions for table columns
func (t *Table) Column(name string) contract.Column {
	c := column.NewColumn(t.name, name, t.connect)
	t.columns = append(t.columns, c)
	return c
}

func (t *Table) String(name string, length int) contract.Column {
	if length < 0 {
		length = 255
	}
	c := t.Column(name).Type("varchar(" + strconv.Itoa(length) + ")")
	return c
}

func (t *Table) Integer(name string) contract.Column {
	c := t.Column(name).Type("int")
	return c
}

func (t *Table) WithTimestamps() contract.Table {
	t.timestamps = true
	return t
}

func (t *Table) getTimeStampsSQL() string {
	return "created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP," +
		"updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP"
}

func (t *Table) RenameColumn(oldName, newName string) contract.Column {
	c := t.Column(oldName).Rename(newName)
	return c
}

func (t *Table) DropColumn(name string) contract.Column {
	c := t.Column(name).Drop()
	return c
}

// Functions for keys
func (t *Table) PrimaryKey(Column string) contract.Column {
	if ok, c := t.hasColumn(Column); ok {
		return c.Primary()
	} else {
		return t.Integer("id").Primary()
	}
}

func (t *Table) ForeignKey(Column string) contract.ForeignKey {
	k := key.NewForeignKey(t.name, Column)
	t.foreignKeys = append(t.foreignKeys, k)
	return k
}

func (t *Table) DropForeignKey(name string) contract.Table {
	if !t.change {
		return t
	}
	ki := info.GetKeyInfo(name, t.connect)
	if ki == nil {
		panic("Foreign key " + name + " not exist. Foreign key name is correct?")
	}
	t.foreignKeys = append(t.foreignKeys, key.NewForeignKeyByKeyInfo(ki).Drop())
	return t
}

// Generate SQL functions
func (t *Table) GetSQL() string {
	if t.drop {
		return t.dropTableSQL()
	}
	if t.newName != "" {
		return t.renameTableSQL()
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
	return strings.TrimRight(sql, ",") + ");"
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
		forDrop = base + strings.TrimRight(forDrop, ",") + ";"
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

func (t *Table) renameTableSQL() string {
	return fmt.Sprintf("ALTER TABLE %v RENAME TO %v;", t.name, t.newName)
}

// Execution functions
func (t *Table) Exec() error {
	queries := strings.Split(t.GetSQL(), ";")
	for i, q := range queries {
		if i == len(queries) {
			break
		}
		if config.GetConfig().Verbose {
			log.Println(q + ";")
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
		if config.GetConfig().Verbose {
			log.Println(q + ";")
		}
		t.connect.MustExec(q + ";")
	}

	for _, k := range t.uniqueKeys {
		k.MustExec(t.connect)
	}
}

// Helpful functions
func (t *Table) checkUniqueKey(c contract.Column) {
	if c.NeedUniqueKey() {
		k := key.NewUniqueKey(t.name, c.GetName()).GenerateKeyName()
		t.uniqueKeys = append(t.uniqueKeys, k)
	}
	if c.NeedDropUniqueKey() {
		k := key.NewUniqueKey(t.name, c.GetName()).SetKeyName(c.GetUniqueKeyName()).Drop()
		t.uniqueKeys = append(t.uniqueKeys, k)
	}
}

// Helpful functions
func (t *Table) hasColumn(name string) (bool, contract.Column) {
	for _, c := range t.columns {
		if c.GetName() == name {
			return true, c
		}
	}
	return false, nil
}
