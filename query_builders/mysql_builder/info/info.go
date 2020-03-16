package info

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ColumnInfo struct {
	ColumnDefault  sql.NullString `db:"COLUMN_DEFAULT"`
	IsNullable     string         `db:"IS_NULLABLE"`
	ColumnType     string         `db:"COLUMN_TYPE"`
	ColumnKey      string         `db:"COLUMN_KEY"`
	ColumnKeysInfo []*IndexInfo
}

func (ci *ColumnInfo) Nullable() bool {
	return ci.IsNullable == "YES"
}

func (ci *ColumnInfo) HasDefault() (bool, string) {
	hasDefault := ci.ColumnDefault.String != ""
	return hasDefault, ci.ColumnDefault.String
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

func GetColumnInfo(table, column string, db *sqlx.DB) *ColumnInfo {
	ci := &ColumnInfo{}
	_ = db.Get(ci, "SELECT COLUMN_DEFAULT, IS_NULLABLE, COLUMN_TYPE, COLUMN_KEY FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name=? AND COLUMN_NAME=?;",
		table, column)
	if ci.ColumnType == "" {
		return nil
	}
	ci.ColumnKeysInfo = GetIndexInfoByColumn(table, column, db)
	return ci
}

type IndexInfo struct {
	NonUnique int    `db:"Non_unique"`
	KeyName   string `db:"Key_name"`
}

func (ii *IndexInfo) IsUnique() bool {
	return ii.NonUnique != 1 && ii.KeyName != "PRIMARY"
}

func (ii *IndexInfo) IsPrimary() bool {
	return ii.NonUnique != 1 && ii.KeyName == "PRIMARY"
}

func (ii *IndexInfo) IsForeign() bool {
	return ii.NonUnique != 0
}

func GetIndexInfoByColumn(table, column string, db *sqlx.DB) []*IndexInfo {
	var ii []*IndexInfo
	_ = db.Select(&ii, fmt.Sprintf("SHOW INDEXES FROM %v WHERE Column_name='%v';", table, column))
	return ii
}

func GetIndexInfoByTable(table string, db *sqlx.DB) []*IndexInfo {
	var ii []*IndexInfo
	_ = db.Select(&ii, fmt.Sprintf("SHOW INDEXES FROM %v;", table))
	return ii
}

type KeyInfo struct {
	ConstraintName       string `db:"CONSTRAINT_NAME"`
	TableName            string `db:"TABLE_NAME"`
	ColumnName           string `db:"COLUMN_NAME"`
	ReferencedTableName  string `db:"REFERENCED_TABLE_NAME"`
	ReferencedColumnName string `db:"REFERENCED_COLUMN_NAME"`
}

func GetKeyInfo(name string, db *sqlx.DB) *KeyInfo {
	var ki []*KeyInfo
	_ = db.Select(&ki, fmt.Sprintf("SELECT CONSTRAINT_NAME, TABLE_NAME, COLUMN_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE CONSTRAINT_NAME='%v';", name))
	if len(ki) < 1 {
		return nil
	}
	return ki[0]
}
