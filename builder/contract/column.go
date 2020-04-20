package contract

type Column interface {
	Type(string) Column
	Nullable() Column
	NotNull() Column
	Autoincrement() Column
	NotAutoincrement() Column
	Default(string) Column
	Primary() Column
	Unique() Column
	NotUnique() Column
	Drop() Column
	Change() Column
	First() Column
	After(string) Column
	Rename(string) Column
	GetSQL() string
	GetName() string
	GetUniqueKeyName() string
	IsPrimary() bool
	IsUnique() bool
	HasUniqueKey() bool
	NeedUniqueKey() bool
	NeedDropUniqueKey() bool
	IsWaitingDrop() bool
	IsWaitingRename() bool
	IsWaitingChange() bool
}
