After installation, let's just hit metana on the terminal.

```shell
$ metana

An abstract migration tool for all types of migrations

Usage:
  metana [command]

Available Commands:
  config      Manage your local metana config in .metana.yml
  create      Create a migration in Go
  down        Run downward migrations
  help        Help about any command
  init        Initialize a migrations directory
  list        List existing migrations
  up          Run upward migrations

Flags:
      --config string   config gen (default is $HOME/.metana.yaml)
  -h, --help            help for metana
  -t, --toggle          Help message for toggle

Use "metana [command] --help" for more information about a command.
```

### `Init`

`init` initializes a boilerplate migrations directory in your current path.

```shell
$ metana init
 ✓ Created /Users/g14a/metana/migrations/main.go
```

By default it will create a `migrations` folder if no such folder exists. If it does, it adds the `main.go` file into the same.

If you want to initialize migrations in a different directory, you can do so with the `--dir | -d` flag:

```shell
metana init --dir /path/to/folder
```

### `Create`

`create` creates a migration script with two functions `Up()` and `Down()` denoting the upward and downward migration of the same.

```shell
$ metana create initSchema
 ✓ Created /Users/g14a/metana/migrations/1614532908-Sample.go
 ✓ Updated /Users/g14a/metana/migrations/main.go
```

Head over to your `1614532908-InitSchema.go` to edit your script. Remember to not change any function signature. 

### `Up`

`up` runs all the upward migrations in the migrations directory in order of their creation time.

```shell
$ metana up

  >>> Migrating up: 1619942687-InitSchema.go
InitSchema up

  >>> Migrating up: 1619942704-AddIndexes.go
AddIndexes up

  >>> migration : complete
```

### `Down`

`down` runs the downward migrations in the reverse order of creation time because we're trying to undo the upward migrations.

```shell
$ metana down

  >>> Migrating down: 1619942704-AddIndexes.go
  AddIndexes down

  >>> Migrating down: 1619942687-InitSchema.go
  InitSchema down

  >>> migration : complete
```

### `List`

`list` lists all the migrations present in your migrations folder along with the last modified time.

```shell
$ metana list

  +----------------------------------+------------------+
  |            MIGRATION             |  LAST MODIFIED   |
  +----------------------------------+------------------+
  | 1619943670-InitSchema.go         | 02-05-2021 13:51 |
  | 1619943677-AddIndexes.go         | 02-05-2021 13:51 |
  | 1619943874-AddFKeys.go           | 02-05-2021 13:54 |
  | 1619943888-AddBoilerPlateRows.go | 02-05-2021 13:54 |
  +----------------------------------+------------------+
```