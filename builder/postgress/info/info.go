package info

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type ColumnInfo struct {
	ColumnDefault   sql.NullString `db:"column_default"`
	IsNullable      string         `db:"is_nullable"`
	ColumnType      string         `db:"data_type"`
	CharacterLength int            `db:"character_maximum_length"`
	ColumnKeysInfo  []*IndexInfo
}

func (ci *ColumnInfo) Nullable() bool {
	return ci.IsNullable == "YES"
}

func (ci *ColumnInfo) HasDefault() (bool, string) {
	if ci.ColumnDefault.String == "" {
		return false, ""
	}
	if strings.Contains(ci.ColumnDefault.String, "'") {
		return true, strings.Split(ci.ColumnDefault.String, "'")[1]
	}
	return true, ci.ColumnDefault.String
}

func (ci *ColumnInfo) IsUnique() bool {
	for _, k := range ci.ColumnKeysInfo {
		if k.IsUnique() {
			return true
		}
	}
	return false
}

func (ci *ColumnInfo) GetUniqueKey() *IndexInfo {
	for _, k := range ci.ColumnKeysInfo {
		if k.IsUnique() {
			return k
		}
	}
	return nil
}

func (ci *ColumnInfo) IsPrimary() bool {
	for _, k := range ci.ColumnKeysInfo {
		if k.IsPrimary() {
			return true
		}
	}
	return false
}

func (ci *ColumnInfo) GetType() string {
	if ci.ColumnType == "character varying" {
		return fmt.Sprintf("varchar(%v)", ci.CharacterLength)
	}
	return ci.ColumnType
}

func GetColumnInfo(table, column string, db *sqlx.DB) *ColumnInfo {
	ci := &ColumnInfo{}
	_ = db.Get(ci, "SELECT column_default, is_nullable, data_type, character_maximum_length FROM information_schema.columns"+
		" WHERE table_name=$1 AND column_name=$2;",
		table, column)
	if ci.ColumnType == "" {
		return nil
	}
	ci.ColumnKeysInfo = GetIndexInfoByColumn(table, column, db)
	return ci
}

type IndexInfo struct {
	KeyName string `db:"constraint_name"`
	KeyType string `db:"constraint_type"`
}

func (ii *IndexInfo) IsUnique() bool {
	return ii.KeyType == "UNIQUE"
}

func (ii *IndexInfo) IsPrimary() bool {
	return ii.KeyType == "PRIMARY KEY"
}

func (ii *IndexInfo) IsForeign() bool {
	return ii.KeyType == "FOREIGN KEY"
}

func GetIndexInfoByColumn(table, column string, db *sqlx.DB) []*IndexInfo {
	var ii []*IndexInfo
	_ = db.Select(&ii, "SELECT ccu.constraint_name, tc.constraint_type FROM information_schema.constraint_column_usage as ccu"+
		" INNER JOIN information_schema.table_constraints as tc ON tc.constraint_name = ccu.constraint_name "+
		" WHERE ccu.table_name = $1 AND ccu.column_name = $2;", table, column)
	return ii
}

func GetIndexInfoByTable(table string, db *sqlx.DB) []*IndexInfo {
	var ii []*IndexInfo
	_ = db.Select(&ii, "SELECT ccu.constraint_name, tc.constraint_type FROM information_schema.constraint_column_usage as ccu"+
		" INNER JOIN information_schema.table_constraints as tc ON tc.constraint_name = ccu.constraint_name "+
		" WHERE ccu.table_name = $1;", table)
	return ii
}

type KeyInfo struct {
	ConstraintName string `db:"constraint_name"`
	BaseTable      string `db:"base_table"`
	BaseColumn     string `db:"base_column"`
	TargetTable    string `db:"target_table"`
	TargetColumn   string `db:"target_column"`
	OnUpdate       string `db:"update_rule"`
	OnDelete       string `db:"delete_rule"`
}

func GetKeyInfo(name string, db *sqlx.DB) *KeyInfo {
	var ki []*KeyInfo
	_ = db.Select(&ki, "SELECT rc.constraint_name,update_rule,delete_rule, kcu.column_name as base_column, kcu.table_name as base_table, ccu.column_name as target_column, ccu.table_name as target_table"+
		" FROM information_schema.referential_constraints as rc"+
		" INNER JOIN information_schema.key_column_usage as kcu ON kcu.constraint_name = rc.constraint_name"+
		" INNER JOIN information_schema.constraint_column_usage as ccu ON ccu.constraint_name = rc.constraint_name"+
		" WHERE rc.constraint_name = $1;", name)
	if len(ki) < 1 {
		return nil
	}
	return ki[0]
}
