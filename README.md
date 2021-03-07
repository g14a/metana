# go-migrate

![demo.gif](https://github.com/g14a/go-migrate/blob/main/demo.gif)

```shell
An abstract migration framework for all types of migrations

Usage:
  go-migrate [command]

Available Commands:
  create      create a migration script in Go
  down        Run the downward migration
  help        Help about any command
  init        initialize a migrations directory
  list        list migrations
  up          Run the upward migration

Flags:
      --config string   config gen (default is $HOME/.go-migrate.yaml)
  -h, --help            help for go-migrate
  -t, --toggle          Help message for toggle

Use "go-migrate [command] --help" for more information about a command.
```

### Steps to create a migration

```shell
# Init migration
> go-migrate init
 ✓ Created /Users/g14a/go-migrate/migrations/main.go

# Create a migration
> go-migrate create sample
 ✓ Created migrations/1614532908-Sample.go
 ✓ Generated migrations/main.go
 
# Run upward migration
> go-migrate up
Sample up

# Run downward migration
> go-migrate down
Sample down

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

