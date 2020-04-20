package contract

import "github.com/jmoiron/sqlx"

type UniqueKey interface {
	SetKeyName(string) UniqueKey
	GenerateKeyName() UniqueKey
	Drop() UniqueKey
	GetSQL() string
	Exec(*sqlx.DB) error
	MustExec(*sqlx.DB)
	GetName() string
}
