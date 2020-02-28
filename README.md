## Overview

Simple module for using migrations in your project
 
`Now only for MySQL!`

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
```

In `migrations/list` directory create your migrations like existed example

In `migrations/entry.go` in `getMigrationsList()` method put your migrations structures