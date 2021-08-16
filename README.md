## Overview

Module for the creating migrations for GoLang applications.
 
`Now only for MySQL and Postgres!`

### Install

```
    go get github.com/ShkrutDenis/go-migrations
```

for update use flag `-u`:

```
    go get -u github.com/ShkrutDenis/go-migrations
```

### Usage

Run this command for put to your project the template for usage go-migrations:
```
    bash $GOPATH/github.com/ShkrutDenis/go-migrations/init.sh
    or if you use vendor folder
    bash vendor/github.com/ShkrutDenis/go-migrations/init.sh
```

Or you can copy sources from your dependencies path manually if you have trouble with command.
For example from:
```
    .../github.com/ShkrutDenis/go-migrations/template
```

In `migrations/list` directory create your migrations like existed example

In `migrations/entry.go` in `getMigrationsList()` method put your migrations structures

For migrate:
```
    go run migrations/entry.go
```

If you want to rollback, add `--rollback` flag.

#### Environment variables

Module uses next variables for creating a connection with DB:

- DB_DRIVER
- DB_USER
- DB_PASSWORD
- DB_HOST
- DB_PORT
- DB_NAME

if `DB_HOST` and `DB_PORT` doesn't exist, will be used a `DB_CONNECTION` with next format: `host:port`

By default, module load env file from the current directory with name `.env`. For use custom env file you can use next flags: `--env-path` and `--env-file`

#### Directly usage

If you need run migrations from your application, then you could use migration engine by next way:
```go
    package your_package

    import (
        gmEngine "github.com/ShkrutDenis/go-migrations/engine"
        gmConfig "github.com/ShkrutDenis/go-migrations/engine/config"
        gmStore "github.com/ShkrutDenis/go-migrations/engine/store"
    )

    func Migrate() {
        e := gmEngine.NewEngine()
        e.WithConfig(gmConfig.Config{
        	Verbose:    true,          // set as true if you want to see additional logs 
        	IsRollback: false,         // set as true if you want to rollback migrations 
        	EnvFile:    ".custom.env", // default value is `.env`. If ypu need to override set value here 
        	EnvPath:    "path",        // default value is empty (current directory). If you need to override set value here
        })
        e.Run(getMigrationsList())
    }
    
    func getMigrationsList() []gmStore.Migratable {
    	return []gmStore.Migratable{
            // Register you migrations here
        }
    }
```

If you want to use your DB connection instead of the created by library, you could use next way:

```go
    func Migrate() {
        e := gmEngine.NewEngine()
        
        e.GetConnector().SetConnection() // if you use sql.DB
        // or
        e.GetConnector().SetConnectionX() // if you use sqlx.DB
        
        e.Run(getMigrationsList())
    }
```

### Documentation

You can found documentation [here](https://github.com/ShkrutDenis/go-migrations/tree/master/docs).

### Examples

You can found few examples with create migrations [here](https://github.com/ShkrutDenis/go-migrations/tree/master/examples).

### License

Licensed under [MIT License](https://github.com/ShkrutDenis/go-migrations/blob/master/LICENSE)
