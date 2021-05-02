# go-migrate

An abstract migration tool for all types of migrations

![demo.gif](https://github.com/g14a/go-migrate/blob/main/demo.gif)

### Install

```shell
go get github.com/g14a/go-migrate
```

### Build from source
Make sure you have Go installed.
```shell
git clone https://github.com/g14a/go-migrate.git
cd go-migrate
go install
```

### Usage
```shell
An abstract migration tool for all types of migrations

Usage:
  go-migrate [command]

Available Commands:
  create      Create a migration in Go
  down        Run downward migrations
  help        Help about any command
  init        Initialize a migrations directory
  list        List existing migrations
  up          Run upward migrations

Flags:
      --config string   config gen (default is $HOME/.go-migrate.yaml)
  -h, --help            help for go-migrate
  -t, --toggle          Help message for toggle

Use "go-migrate [command] --help" for more information about a command.```
```

### Steps to create a migration

```shell
# Init migration
> go-migrate init
 ✓ Created /Users/g14a/go-migrate/migrations/main.go
 ✓ Created /Users/g14a/go-migrate/migrations/store.go
 ✓ Created /Users/g14a/go-migrate/migrations/migrate.json

# Create a migration
> go-migrate create sample
 ✓ Created /Users/g14a/go-migrate/migrations/1614532908-Sample.go
 ✓ Generated /Users/g14a/go-migrate/migrations/main.go
 
# Run upward migration
> go-migrate up

  >>> Migrating up: 1614532908-Sample.go
  Sample up

  >>> migration : complete

# Run downward migration
> go-migrate down
  
  >>> Migrating down: 1614532908-Sample.go
  Sample down

  >>> migration : complete

# List migrations
> go-migrate list
+----------------------+------------------+
|      MIGRATION       |  LAST MODIFIED   |
+----------------------+------------------+
| 1614532908-Sample.go | 28-02-2021 22:51 |
+----------------------+------------------+

Navigate to the file 1614532908-Sample.go inside
the migrations directory to update your migration script.
```

### Features

* Specify custom migrations directory when initializing, creating and running migrations with ```--dir | -d ``` flag. By default ```dir``` is set to ```migrations```

```shell
> go-migrate init --dir schema-mig
 ✓ Created /Users/g14a/go-migrate/schema-mig/main.go
 ✓ Created /Users/g14a/go-migrate/schema-mig/store.go
 ✓ Created /Users/g14a/go-migrate/schema-mig/migrate.json

> go-migrate create initSchema --dir schema-mig
 ✓ Created /Users/g14a/go-migrate/schema-mig/scripts/1619943164-InitSchema.go
 ✓ Generated /Users/g14a/go-migrate/schema-mig/main.go
 
> go-migrate up --dir schema-mig
  >>> Migrating up: 1619943670-InitSchema.go
InitSchema up

  >>> migration : complete
```

* Run upward and downward migrations until(including) a certain migration with the ```--until``` flag.

```shell
> go-migrate init                                                                              
 ✓ Created /Users/g14a/go-migrate/migrations/main.go
 ✓ Created /Users/g14a/go-migrate/migrations/store.go
 ✓ Created /Users/g14a/go-migrate/migrations/migrate.json
 
> go-migrate create initSchema                                                                
 ✓ Created /Users/g14a/go-migrate/migrations/scripts/1619942687-InitSchema.go
 ✓ Generated /Users/g14a/go-migrate/migrations/main.go
 
> Create more migration scripts...

> go-migrate list

  +----------------------------------+------------------+
  |            MIGRATION             |  LAST MODIFIED   |
  +----------------------------------+------------------+
  | 1619943670-InitSchema.go         | 02-05-2021 13:51 |
  | 1619943677-AddIndexes.go         | 02-05-2021 13:51 |
  | 1619943874-AddFKeys.go           | 02-05-2021 13:54 |
  | 1619943888-AddBoilerPlateRows.go | 02-05-2021 13:54 |
  +----------------------------------+------------------+

> go-migrate up --until AddFkeys                                                                

  >>> Migrating up: 1619942687-InitSchema.go
InitSchema up

  >>> Migrating up: 1619942704-AddIndexes.go
AddIndexes up

  >>> Migrated up until: 1619942704-AddIndexes.go

  >>> migration : complete

> go-migrate down --until AddIndexes
  
  >>> Migrating down: 1619943888-AddBoilerPlateRows.go
AddBoilerPlateRows down

  >>> Migrating down: 1619943874-AddFKeys.go
AddFKeys down

  >>> Migrating down: 1619943677-AddIndexes.go
AddIndexes down

  >>> Migrated down until: 1619943677-AddIndexes.go

  >>> migration : complete
```

