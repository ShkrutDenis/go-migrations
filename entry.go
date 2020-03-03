package main

import (
	"flag"
	gm "github.com/ShkrutDenis/go-migrations"
	gmStore "github.com/ShkrutDenis/go-migrations/store"
	"go-migrations-local/list"

	//"github.com/ShkrutDenis/go-migrations/template/migrations/list"
)

var isRollback *bool

func init() {
	isRollback = flag.Bool("rollback", false, "")
	flag.Parse()
}

func main() {
	if *isRollback {
		gm.Rollback(getMigrationsList())
		return
	}

	gm.Migrate(getMigrationsList())
}

func getMigrationsList() []gmStore.Migratable {
	return []gmStore.Migratable{
		&list.CreateExample1Table{},
		&list.CreateExample2Table{},
		&list.UpdateExample1Table{},
		// Register you migrations here
	}
}
