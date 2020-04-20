package contract

type ForeignKey interface {
	Reference(string) ForeignKey
	On(string) ForeignKey
	OnUpdate(string) ForeignKey
	OnDelete(string) ForeignKey
	Drop() ForeignKey
	SetKeyName(string) ForeignKey
	GenerateKeyName() ForeignKey
	GetSQL() string
	GetName() string
	ForDrop() bool
}
