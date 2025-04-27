![](metana.png)

# Metana

![OpenSource](https://img.shields.io/badge/Open%20Source-000000?style=for-the-badge&logo=github)
![go](https://img.shields.io/badge/-Written%20In%20Go-00add8?style=for-the-badge&logo=Go&logoColor=ffffff)
![cli](https://img.shields.io/badge/-Build%20for%20CLI-000000?style=for-the-badge&logo=Powershell&logoColor=ffffff)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/g14a/metana.svg)](https://github.com/g14a/metana)
[![Go Report Card](https://goreportcard.com/badge/github.com/g14a/metana)](https://goreportcard.com/report/github.com/g14a/metana)
[![Go Workflow Status](https://github.com/g14a/metana/workflows/Go/badge.svg)](https://github.com/g14a/metana/workflows/Go/badge.svg)

An abstract task migration tool written in Go for Go services. Database and non database migrations management brought to your CLI.

# Table of Contents

* [Use case](https://github.com/g14a/metana#use-case)
* [Installation](https://github.com/g14a/metana#installation)
    * [Using Go](https://github.com/g14a/metana#using-go)
    * [macOS](https://github.com/g14a/metana#mac)
    * [Linux](https://github.com/g14a/metana#linux)
    * [Building from Source](https://github.com/g14a/metana#building-from-source)
    * [Docker](https://github.com/g14a/metana#docker)
* [Usage](https://github.com/g14a/metana#usage)
    * [Init](https://github.com/g14a/metana#init) ✅
    * [Create](https://github.com/g14a/metana#create) 👌 
    * [Up](https://github.com/g14a/metana#up) ⬆️
    * [Down](https://github.com/g14a/metana#down) ⬇️
    * [List](https://github.com/g14a/metana#list) 
* [Features](https://github.com/g14a/metana#features)
    * [Custom directory to store migrations](https://github.com/g14a/metana#custom-directory-to-store-migrations)
    * [Run migrations until a certain point](https://github.com/g14a/metana#run-a-migration-until-a-certain-point)
    * [Store and Track your migrations in your favourite database](https://github.com/g14a/metana#store-and-track-your-migrations-in-your-favourite-database)
    * [Dry Run Migrations](https://github.com/g14a/metana#dry-run-migrations)
    * [Custom Config](https://github.com/g14a/metana#custom-config)
    * [Automatic Rollback on Migration Failure](https://github.com/g14a/metana#automatic-rollback-on-migration-failure)
    
# Use case

The motivation behind creating this tool, is to abstract away the database part. If your task can be completed with Pure Go or via a Go driver of your service, then this is for you. Since it makes use of the Go runtime, you can even perform database migrations like PostgreSQL, Mongo, Redis, Elasticsearch, GCP Buckets etc. You just need to be able to interact with your data store or complete your task using Go.

The main use case is when you won't be able to do everything with SQL or No-SQL syntax. There might be some tasks where you need to aggregate data, iterate over them, and do business related stuff with the retrieved data. All you need to know is Go syntax and write a Go program.

# Installation

## Using Go
```shell
go get github.com/g14a/metana
```

## **Mac**

```shell
brew tap g14a/homebrew-metana
brew install metana
```
## **Linux**

Checkout the releases page and download your platform's binaries to install them.

[Releases Page](https://github.com/g14a/metana/releases)

## **Building from source**

Prerequisites:

* Git
* Go 1.13 or newer. Go modules are needed. Better if its the latest version.

```shell
git clone https://github.com/g14a/metana
cd metana
go install
```

# Usage

After installation, let's just hit metana on the terminal.

```shell
$ metana
An abstract migration tool for Go services

Usage:
  metana [flags]
  metana [command]

Available Commands:
  completion  Generate shell completion script
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

## **`Init`**

`init` initializes a boilerplate migrations directory in your current path.

```shell
$ metana init
Successfully initialized migration setup in migrations
```

By default it will create a `migrations` folder if no such folder exists. If it does, it adds the `main.go` file into the same.

If you want to initialize migrations in a different directory, you can do so with the `--dir | -d` flag:

```shell
metana init --dir /path/to/folder
```

## **`Create`**

`create` creates a migration script with two functions `Up()` and `Down()` denoting the upward and downward migration of the same.

```shell
$ metana create initSchema
 ✓ .metana.yml found
 ✓ Created /Users/gowtham.munukutla/metana/migrations/scripts/1745742878_initSchema.go
```

Head over to your `1745742878_initSchema.go` to edit your script. Remember to not change any function signature.

## **`Up`**

`up` runs all the upward migrations in the migrations directory in order of their creation time.

```shell
$ metana up
 ✓ .metana.yml found
InitSchema up
__COMPLETE__[up]: 1745742878_initSchema.go
InitSchema2 up
__COMPLETE__[up]: 1745742917_initSchema2.go
  >>> migration : complete
```

## **`Down`**

`down` runs the downward migrations in the reverse order of creation time because we're trying to undo the upward migrations.

```shell
$ metana down
 ✓ .metana.yml found
InitSchema down
__COMPLETE__[down]: 1745742878_initSchema.go
InitSchema2 down
__COMPLETE__[down]: 1745742917_initSchema2.go
  >>> migration : complete
```

## **`List`**

`list` lists all the migrations present in your migrations folder along with the last executed time.

```shell
$ metana list
+---------------------------+------------------+
|         MIGRATION         |   EXECUTED AT    |
+---------------------------+------------------+
| 1745742878_initSchema.go  | 27-04-2025 14:06 |
| 1745742917_initSchema2.go | 27-04-2025 14:06 |
+---------------------------+------------------+
```

# Features

## **Custom directory to store migrations**

Specify a custom directory when creating and running upward or downward migrations using the `--dir` flag. Be default it is set to `"migrations"`

```shell
$ metana init --dir custom-migration-directory
Successfully initialized migration setup in custom-migration-directory

$ metana create initSchema --dir custom-migration-directory
 ✓ Created /Users/gowtham.munukutla/metana/custom-migration-directory/scripts/1745743111_initSchema.go
 
$ metana up --dir custom-migration-directory
InitSchema up
__COMPLETE__[up]: 1745743111_initSchema.go
  >>> migration : complete
```

## **Run a migration until a certain point**

Run upward and downward migrations until(and including) a certain migration with the `--until` flag.

```shell

$ metana create initSchema                                                  ✓ .metana.yml found
 ✓ Created /Users/gowtham.munukutla/metana/migrations/scripts/1745743242_initSchema.go
 
$ Create more migration scripts...

$ metana list

+---------------------------+------------------+
|         MIGRATION         |   EXECUTED AT    |
+---------------------------+------------------+
| 1745743242_initSchema.go  |                  |
| 1745743245_initSchema2.go |                  |
| 1745743247_initSchema3.go |                  |
+---------------------------+------------------+

$ metana up --until initSchema2                                                 ✓ .metana.yml found
InitSchema up
__COMPLETE__[up]: 1745743242_initSchema.go
InitSchema2 up
__COMPLETE__[up]: 1745743245_initSchema2.go
 >>> Reached --until: initSchema2. Stopping further migrations.
  >>> migration : complete
```

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

You can dry run your migrations using the explicit `--dry` option. This option doesn't track any migrations, doesn't create a default `migrate.json` file. It literally just dry runs. However your tasks are run. This helps when you're incrementally writing, testing and running your functions instead of manually deleting states in your store.

```shell
$ metana up --dry
 ✓ .metana.yml found
InitSchema up
__COMPLETE__[up]: 1745743242_initSchema.go
InitSchema2 up
__COMPLETE__[up]: 1745743245_initSchema2.go
InitSchema3 up
__COMPLETE__[up]: 1745743247_initSchema3.go
  >>> dry run migration : complete
```

```shell
$ metana down --dry
 ✓ .metana.yml found
InitSchema down
__COMPLETE__[down]: 1745743242_initSchema.go
InitSchema2 down
__COMPLETE__[down]: 1745743245_initSchema2.go
InitSchema3 down
__COMPLETE__[down]: 1745743247_initSchema3.go
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

## **Automatic Rollback on Migration Failure**

Metana automatically handles rollback during upward migrations (`metana up`)

If an **upward migration** (`up`) fails while being run, Metana will immediately **trigger the downward migration** (`down`) of that **same migration file** to rollback the changes and restore consistency.

You don't have to manually clean up — rollback is automatic. But you still have implement the logic of the downward migration.

**Example:**

```shell
$ metana up
 ✓ .metana.yml found
InitSchema up
__COMPLETE__[up]: 1745748076_initSchema.go
Migration 1745748078_initSchema2.go failed, attempting rollback...
InitSchema2 down
__COMPLETE__[down]: 1745748078_initSchema2.go
  >>> migration : complete
2025/04/27 15:32:05 migration 1745748078_initSchema2.go failed: execution error: exit status 1
error: simulated error
goroutine 1 [running]:
runtime/debug.Stack()
        /Users/gowtham.munukutla/.gvm/gos/go1.21/src/runtime/debug/stack.go:24 +0x64
runtime/debug.PrintStack()
        /Users/gowtham.munukutla/.gvm/gos/go1.21/src/runtime/debug/stack.go:16 +0x1c
main.main()
        /Users/gowtham.munukutla/metana/migrations/scripts/1745748078_initSchema2.go:48 +0x1d8
exit status 1
```