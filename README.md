## Overview

A simple module for using migrations in your project
 
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
    bash $GOPATH/src/github.com/ShkrutDenis/go-migrations/init.sh
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

if `DB_HOST` and `DB_PORT` doesn’t exist, will be used a `DB_CONNECTION` with next format: `host:port`

By default, module load env file from the current directory with name `.env`. For use custom env file you can use next flags: `--env-path` and `--env-file`

#### Examples

You can found few examples with create migration in `/examples` repository folder.