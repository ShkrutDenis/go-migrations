package mysql_builder

const (
	nullable      = "null"
	notNullable   = "not null"
	autoIncrement = "auto_increment"
)

type Field struct {
	name          string
	fieldType     string
	hasDefault    bool
	defaultValue  string
	nullable      string
	autoincrement bool
}

func NewField(fieldType, fieldName string) *Field {
	return &Field{name: fieldName, fieldType: fieldType, nullable: notNullable}
}

func (f *Field) Nullable() *Field {
	f.nullable = nullable
	return f
}

func (f *Field) Autoincrement() *Field {
	f.autoincrement = true
	return f
}

func (f *Field) Default(value string) *Field {
	f.hasDefault = true
	f.defaultValue = value
	return f
}

func (f *Field) GetSQL() string {
	sql := f.name + " " + f.fieldType
	if f.hasDefault {
		sql += " " + f.defaultValue
	}
	if f.autoincrement {
		sql += " " + autoIncrement
	} else {
		sql += " " + f.nullable
	}
	return sql
}
