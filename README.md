# metana

![OpenSource](https://img.shields.io/badge/Open%20Source-000000?style=for-the-badge&logo=github)
![go](https://img.shields.io/badge/-Written%20In%20Go-00add8?style=for-the-badge&logo=Go&logoColor=ffffff)
![cli](https://img.shields.io/badge/-Build%20for%20CLI-000000?style=for-the-badge&logo=Powershell&logoColor=ffffff)

[![Go Report Card](https://goreportcard.com/badge/github.com/g14a/metana)](https://goreportcard.com/report/github.com/g14a/metana)
[![Go Workflow Status](https://github.com/g14a/metana/workflows/Go/badge.svg)](https://github.com/g14a/metana/workflows/Go/badge.svg)

An abstract task migration tool written in Go for Go services. Database and non database migrations management brought to your CLI.

![demo.gif](https://github.com/g14a/metana/blob/main/demo.gif)

## Use case
The motivation behind creating this tool, is to abstract away the database part. If your task can be completed with Pure Go or via a Go driver of 
your service, then this is for you. Since it makes use of the Go runtime, you can even perform database migrations like PostgreSQL, Mongo, Redis, Elasticsearch, GCP Buckets etc.
You just need to be able to interact with your data store or complete your task using Go.

The main use case is when you won't be able to do everything with SQL or No-SQL syntax. 
There might be some tasks where you need to aggregate data, iterate over them, and do business related stuff with the retrieved data.
All you need to know is Go syntax and write a Go program.
### Install

```shell
go get github.com/g14a/metana
```

Checkout the [release binaries](https://github.com/g14a/metana/releases) page for precompiled binaries.

### Build from source
Make sure you have Go installed.
```shell
git clone https://github.com/g14a/metana.git
cd metana
go install
```

### Docker
```shell
docker pull g14a/metana
docker run 
```

### Usage
```shell
An abstract migration tool for all types of migrations

Usage:
  metana [command]

Available Commands:
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

### Steps to create a migration

```shell
# Init migration
> metana init
 ✓ Created /Users/g14a/metana/migrations/main.go

# Create a migration
> metana create sample
 ✓ Created /Users/g14a/metana/migrations/1614532908-Sample.go
 ✓ Generated /Users/g14a/metana/migrations/main.go
 
# Run upward migration
> metana up

  >>> Migrating up: 1614532908-Sample.go
  Sample up

  >>> migration : complete

# Run downward migration
> metana down
  
  >>> Migrating down: 1614532908-Sample.go
  Sample down

  >>> migration : complete

# List migrations
> metana list
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
> metana init --dir schema-mig
 ✓ Created /Users/g14a/metana/schema-mig/main.go

> metana create initSchema --dir schema-mig
 ✓ Created /Users/g14a/metana/schema-mig/scripts/1619943164-InitSchema.go
 ✓ Generated /Users/g14a/metana/schema-mig/main.go
 
> metana up --dir schema-mig
  >>> Migrating up: 1619943670-InitSchema.go
InitSchema up

  >>> migration : complete
```

* Run upward and downward migrations until(including) a certain migration with the ```--until``` flag.

```shell
> metana init                                                                              
 ✓ Created /Users/g14a/metana/migrations/main.go

> metana create initSchema                                                                
 ✓ Created /Users/g14a/metana/migrations/scripts/1619942687-InitSchema.go
 ✓ Generated /Users/g14a/metana/migrations/main.go
 
> Create more migration scripts...

> metana list

  +----------------------------------+------------------+
  |            MIGRATION             |  LAST MODIFIED   |
  +----------------------------------+------------------+
  | 1619943670-InitSchema.go         | 02-05-2021 13:51 |
  | 1619943677-AddIndexes.go         | 02-05-2021 13:51 |
  | 1619943874-AddFKeys.go           | 02-05-2021 13:54 |
  | 1619943888-AddBoilerPlateRows.go | 02-05-2021 13:54 |
  +----------------------------------+------------------+

> metana up --until AddFkeys                                                                

  >>> Migrating up: 1619942687-InitSchema.go
InitSchema up

  >>> Migrating up: 1619942704-AddIndexes.go
AddIndexes up

  >>> Migrated up until: 1619942704-AddIndexes.go

  >>> migration : complete

> metana down --until AddIndexes
  
  >>> Migrating down: 1619943888-AddBoilerPlateRows.go
AddBoilerPlateRows down

  >>> Migrating down: 1619943874-AddFKeys.go
AddFKeys down

  >>> Migrating down: 1619943677-AddIndexes.go
AddIndexes down

  >>> Migrated down until: 1619943677-AddIndexes.go

  >>> migration : complete
```

### Track your migrations in your favourite database

```shell
metana up --store <db-connection-url>
```
```--store``` needs to be passed for downward migrations as well for it to take effect.

Defaults to a ``migrate.json`` if no `url` is provided.

Databases supported for now:
* PostgreSQL
* MongoDB

### Dry run

You can dry run your migrations using the explicit `--dry` option. This option doesn't track any migrations, doesn't create a default `migrate.json` file. It literally just dry runs. However your tasks are run.

```shell
metana up --dry

  >>> Migrating up: 1620656150-Abc.go
Abc up

  >>> Migrating up: 1620656165-Random.go
Random up

  >>> dry run migration : complete
```

```shell
metana down --dry

  >>> Migrating up: 1620656165-Random.go
Random down

  >>> Migrating up: 1620656150-Abc.go
Abc down

  >>> dry run migration : complete
```

All the other options like `--dir` and `--until` work along with `--dry`.

### Set your custom config

Set your custom config in your `.metana.yml` file. As of now it supports `dir` and `store` keys.

For eg:

```
dir: schema-mig
store: '@MONGO_URL'
```
Remember to add it to your git  unless you want to miss migrations on deployments.

If your store has a remote database URL you can specify it via `'@<url>'` syntax and it will automatically be picked up from your environment variables (Remember the single quotes).You don't want to hardcode API Keys and connection URLs in your codebase.

You can either manually add the config on to the `.metana.yml` file or do it via

`metana config set --store @MONGO_URL`

```
$ metana config set --help
Set your metana config

Usage:
  metana config set [flags]

Flags:
  -d, --dir string     Set your migrations directory (default "migrations")
  -h, --help           help for set
  -s, --store string   Set your store
```

CAUTION: If you change the `dir` flag in your `.metana.yml` after running `metana init`, don't forget to rename your migrations directory to the new directory. Otherwise running migrations would result in failure.

### Roadmap
- [ ] Custom Templates