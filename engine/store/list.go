package store

import "github.com/jmoiron/sqlx"

type Migratable interface {
	GetName() string
	Up(*sqlx.DB)
	Down(*sqlx.DB)
}

type Migrations []Migratable

func (m *Migrations) Register(migs []Migratable) {
	*m = migs
}

func (m *Migrations) GetNames() []string {
	var names []string
	for _, i := range *m {
		names = append(names, i.GetName())
	}
	return names
}

func (m *Migrations) Find(name string) Migratable {
	for _, i := range *m {
		if i.GetName() == name {
			return i
		}
	}
	return nil
}

func (m *Migrations) Filter(name string) {
	if (*m)[len(*m)-1].GetName() == name {
		*m = Migrations{}
	}
	for index, i := range *m {
		if i.GetName() == name {
			*m = (*m)[index+1:]
			return
		}
	}
}
