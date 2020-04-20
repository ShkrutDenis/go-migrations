package contract

type Table interface {
	Column(string) Column
	String(string, int) Column
	Integer(string) Column
	WithTimestamps() Table
	RenameColumn(string, string) Column
	DropColumn(string) Column
	PrimaryKey(string) Column
	ForeignKey(string) ForeignKey
	DropForeignKey(string) Table
	GetSQL() string
	Exec() error
	MustExec()
}
