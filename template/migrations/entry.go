package main

import (
	gmEngine "github.com/ShkrutDenis/go-migrations/engine"
	gmStore "github.com/ShkrutDenis/go-migrations/engine/store"
	"github.com/ShkrutDenis/go-migrations/template/migrations/list"
)

func main() {
	e := gmEngine.NewEngine()
	e.Run(getMigrationsList())
	e.GetConnector().Close()
}

func getMigrationsList() []gmStore.Migratable {
	return []gmStore.Migratable{
		&list.CreateExampleTable{},
		// Register you migrations here
	}
}
