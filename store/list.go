package store

import "github.com/jmoiron/sqlx"

type Migratable interface {
	GetName() string
	Up(*sqlx.DB)
	Down(*sqlx.DB)
}

var list []Migratable

func RegisterMigrations(migs []Migratable) {
	list = migs
}

func GetMigrationsList() []Migratable {
	return list
}

func GetMigrationsNames() []string {
	var names []string
	for _, m := range list {
		names = append(names, m.GetName())
	}
	return names
}

func FindMigration(name string) Migratable {
	for _, m := range list {
		if m.GetName() == name {
			return m
		}
	}
	return nil
}

func FilterMigrations(name string) {
	if list[len(list)-1].GetName() == name {
		list = []Migratable{}
	}
	for i, m := range list {
		if m.GetName() == name {
			list = list[i+1:]
		}
	}
}
