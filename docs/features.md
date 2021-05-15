## **Custom directory to store migrations**

Specify a custom directory when creating and running upward or downward migrations using the `--dir` flag. Be default it is set to `"migrations"`

```shell
$ metana init --dir schema-mig
 ✓ Created /Users/g14a/metana/schema-mig/main.go

$ metana create initSchema --dir schema-mig
 ✓ Created /Users/g14a/metana/schema-mig/scripts/1619943164-InitSchema.go
 ✓ Generated /Users/g14a/metana/schema-mig/main.go
 
$ metana up --dir schema-mig
  >>> Migrating up: 1619943670-InitSchema.go
InitSchema up

  >>> migration : complete

```

## **Run a migration until a certain point**

Run upward and downward migrations until(and including) a certain migration with the `--until` flag.

```shell

$ metana create initSchema                                                                
 ✓ Created /Users/g14a/metana/migrations/scripts/1619942687-InitSchema.go
 ✓ Generated /Users/g14a/metana/migrations/main.go
 
$ Create more migration scripts...

$ metana list

  +----------------------------------+------------------+
  |            MIGRATION             |  LAST MODIFIED   |
  +----------------------------------+------------------+
  | 1619943670-InitSchema.go         | 02-05-2021 13:51 |
  | 1619943677-AddIndexes.go         | 02-05-2021 13:51 |
  | 1619943874-AddFKeys.go           | 02-05-2021 13:54 |
  | 1619943888-AddBoilerPlateRows.go | 02-05-2021 13:54 |
  +----------------------------------+------------------+

$ metana up --until AddFkeys                                                                

  >>> Migrating up: 1619942687-InitSchema.go
InitSchema up

  >>> Migrating up: 1619942704-AddIndexes.go
AddIndexes up

  >>> Migrated up until: 1619942704-AddIndexes.go

  >>> migration : complete

$ metana down --until AddIndexes
  
  >>> Migrating down: 1619943888-AddBoilerPlateRows.go
AddBoilerPlateRows down

  >>> Migrating down: 1619943874-AddFKeys.go
AddFKeys down

  >>> Migrating down: 1619943677-AddIndexes.go
AddIndexes down

  >>> Migrated down until: 1619943677-AddIndexes.go

  >>> migration : complete

```

Notice the capitalized format when passing to `--until`.

## **Store and track your migrations in your favourite database**

Store and track your migrations in your favourite database by passing the `--store` flag.

```shell
metana up --store <db-connection-url>
```

If your connection URL is store in an environment variable you can pass it as `--store @MONGO_URL` and it will automatically be picked up from your environment.

Right now, PostgreSQL(which means even CockroachDB URLs) and MongoDB are supported to store migrations.

If no `--store` flag is passed, migrations will be stored in a default `migrate.json` file in the migrations directory.

## **Dry run migrations**

Dry run your migrations using the `--dry` flag.

You can dry run your migrations using the explicit --dry option. This option doesn't track any migrations, doesn't create a default `migrate.json` file. It literally just dry runs. However your tasks are run. This helps when you're incrementally writing, testing and running your functions instead of manually deleting states in your store.

```shell
$ metana up --dry

  >>> Migrating up: 1619942687-InitSchema.go
InitSchema up

  >>> Migrating up: 1619942704-AddIndexes.go
AddIndexes up

  >>> dry run migration : complete
```

```shell
$ metana down --dry

  >>> Migrating down: 1619942704-AddIndexes.go
AddIndexes down

  >>> Migrating down: 1619942687-InitSchema.go
InitSchema down

  >>> dry run migration : complete
```

All the other options like `--dir` and `--until` work along with `--dry`.

## **Custom config**

Set your custom config in your `.metana.yml` file. As of now it supports `dir` and `store` keys.

For eg:
```
dir: schema-mig
store: '@MONGO_URL'
```

Remember to add it to your git unless you want to miss migrations on deployments.

If your store has a remote database URL you can specify it via '@<url>' syntax and it will automatically be picked up from your environment variables (Remember the single quotes).You don't want to hardcode API Keys and connection URLs in your codebase.

`.metana.yml` is created automatically when you run `metana init` which can be used for subsequent migration operations.

You can either manually add the config on to the `.metana.yml` file or do it via

`metana config set --store @MONGO_URL`

```shell
$ metana config set --help
Set your metana config

Usage:
  metana config set [flags]

Flags:
  -d, --dir string     Set your migrations directory (default "migrations")
  -h, --help           help for set
  -s, --store string   Set your store

```

<span style="color:red">CAUTION: </span>
If you change the dir flag in your `.metana.yml` after running `metana init`, don't forget to rename your migrations directory to the new directory. Otherwise running migrations would result in failure.

Priority order of config:

1. Flags passed explicitly
2. `.metana.yml` if it exists.
3. Default values of flags.


## **Wipe out stale migrations**

Wipe out stale(already executed) migration files and update your store with the `wipe` command.

```shell
$ metana wipe

Wipe out old stale migration files and track in your store

Usage:
  metana wipe [flags]

Flags:
  -d, --dir string     Specify custom migrations directory
  -h, --help           help for wipe
  -s, --store string   Specify a connection url to track migrations

Global Flags:
      --config string   config gen (default is $HOME/.metana.yaml)
```

Even the `wipe` command takes configuration from your `.metana.yml` file one exists.
Otherwise the priority order is considered while wiping.