package main

import (
	gm "github.com/ShkrutDenis/go-migrations"
	gmStore "github.com/ShkrutDenis/go-migrations/store"
	"github.com/ShkrutDenis/go-migrations/template/migrations/list"
)

func main() {
	gm.Run(getMigrationsList())
}

func getMigrationsList() []gmStore.Migratable {
	return []gmStore.Migratable{
		&list.CreateExampleTable{},
		// Register you migrations here
	}
}
