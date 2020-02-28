package mysql_builder

import "strconv"

type Table struct {
	name       string
	primaryKey string
	timestamps bool
	fields     []*Field
}

func NewTable(name string) *Table {
	return &Table{name: name}
}

func (t *Table) String(name string, length int) *Field {
	if length < 0 {
		length = 255
	}
	f := NewField(name, "varchar("+strconv.Itoa(length)+")")
	t.fields = append(t.fields, f)
	return f
}

func (t *Table) PrimaryKey(field string) *Table {
	t.primaryKey = field
	return t
}

func (t *Table) WithTimestamps() *Table {
	t.timestamps = true
	return t
}

func (t *Table) GetQuery() string {
	sql := "CREATE TABLE " + t.name + "("
	for i, f := range t.fields {
		sql += f.GetSQL()
		if i != len(t.fields)-1 {
			sql += ","
		}
	}
	if t.timestamps {
		sql += "," + t.getTimeStampsSQL()
	}
	return sql + ");"
}

func (t *Table) getTimeStampsSQL() string {
	return "created_at datetime default current_timestamp not null," +
		"updated_at datetime default current_timestamp ON UPDATE current_timestamp not null"
}
