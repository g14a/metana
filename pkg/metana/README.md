# Metana Migration Library

This package provides a programmatic interface to the metana migration tool, allowing you to run migrations directly from your Go code without using the CLI.

## Installation

```bash
go get github.com/g14a/metana/pkg/metana
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/g14a/metana/pkg/metana"
)

func main() {
    // Create a new metana client
    client := metana.NewClient(metana.Config{
        WorkingDirectory:      ".",
        MigrationsDirectory:   "migrations",
        StoreConnectionString: "postgres://user:pass@localhost/mydb",
    })

    // Initialize migration setup
    if result := client.Init(); !result.Success {
        log.Fatal(result.Error)
    }

    // Create a new migration
    createResult := client.Create(metana.CreateOptions{
        Name: "add_users_table",
    })
    if !createResult.Success {
        log.Fatal(createResult.Error)
    }

    // Run migrations
    upResult := client.Up(metana.UpOptions{})
    if !upResult.Success {
        log.Fatal(upResult.Error)
    }
    
    fmt.Println("Migration completed successfully!")
}
```

## Configuration

The `Config` struct supports the following options:

- **WorkingDirectory**: Base directory for migration operations (defaults to current directory)
- **MigrationsDirectory**: Directory containing migration files (defaults to "migrations")
- **StoreConnectionString**: Connection string for tracking migrations
  - Postgres: `postgres://user:pass@host:port/database`
  - MongoDB: `mongodb://user:pass@host:port/database`
  - File-based (default): Leave empty or provide file path
  - Environment variables: Use `@ENV_VAR_NAME` syntax
- **FileSystem**: Custom filesystem implementation (optional, uses OS filesystem by default)

## API Reference

### Client Creation

```go
client := metana.NewClient(metana.Config{
    // configuration options
})
```

### Operations

#### Initialize Migration Setup

```go
result := client.Init()
```

Creates the migrations directory structure.

#### Create Migration

```go
result := client.Create(metana.CreateOptions{
    Name: "migration_name",
})
```

Creates a new migration file with the specified name.

#### Run Upward Migrations

```go
result := client.Up(metana.UpOptions{
    Until:  "migration_name", // optional: stop at specific migration
    DryRun: false,           // optional: preview without executing
})
```

Executes pending upward migrations.

#### Run Downward Migrations

```go
result := client.Down(metana.DownOptions{
    Until:  "migration_name", // optional: rollback to specific migration
    DryRun: false,           // optional: preview without executing
})
```

Executes downward migrations (rollbacks).

#### List Migrations

```go
result := client.List(metana.ListOptions{})
```

Lists all migrations and their execution status.

### Result Handling

All operations return a `MigrationResult` struct:

```go
type MigrationResult struct {
    Success            bool              // Whether operation succeeded
    Message            string            // Human-readable result message
    Output             string            // Detailed output from operation
    Error              error             // Any error that occurred
    ExecutedMigrations []types.Migration // List of executed migrations
}
```

Always check the `Success` field:

```go
result := client.Up(metana.UpOptions{})
if !result.Success {
    log.Printf("Migration failed: %v", result.Error)
    return
}
fmt.Println("Success:", result.Message)
```

## Store Backends

### File-based Storage (Default)

```go
client := metana.NewClient(metana.Config{
    // StoreConnectionString can be empty for file-based storage
})
```

Stores migration tracking information in `migrations/migrate.json`.

### PostgreSQL

```go
client := metana.NewClient(metana.Config{
    StoreConnectionString: "postgres://user:password@localhost:5432/database",
})
```

### MongoDB

```go
client := metana.NewClient(metana.Config{
    StoreConnectionString: "mongodb://user:password@localhost:27017/database",
})
```

### Environment Variables

```go
client := metana.NewClient(metana.Config{
    StoreConnectionString: "@DATABASE_URL", // Reads from DATABASE_URL env var
})
```

## Testing

You can use an in-memory filesystem for testing:

```go
import "github.com/spf13/afero"

client := metana.NewClient(metana.Config{
    WorkingDirectory:    "/tmp/test",
    MigrationsDirectory: "migrations",
    FileSystem:          afero.NewMemMapFs(),
})
```

## Error Handling

The library provides detailed error information:

```go
result := client.Up(metana.UpOptions{})
if !result.Success {
    // Check for specific error
    if result.Error != nil {
        log.Printf("Error: %v", result.Error)
    }
    
    // Get detailed output
    log.Printf("Output: %s", result.Output)
    
    // Get human-readable message
    log.Printf("Message: %s", result.Message)
}
```

## Migration File Format

Migration files should implement the `Migrator` interface:

```go
package main

import (
    "github.com/g14a/metana/pkg/types"
)

type MyMigration struct{}

func (m MyMigration) Up() error {
    // Your upward migration logic
    return nil
}

func (m MyMigration) Down() error {
    // Your downward migration logic  
    return nil
}

func main() {
    // Migration entry point
    migration := MyMigration{}
    
    // Run based on command line argument
    if len(os.Args) > 1 && os.Args[1] == "-mode" && len(os.Args) > 2 {
        mode := os.Args[2]
        if mode == "up" {
            if err := migration.Up(); err != nil {
                log.Fatal(err)
            }
            fmt.Println("__COMPLETE__[up]:" + "migration_file_name.go")
        } else if mode == "down" {
            if err := migration.Down(); err != nil {
                log.Fatal(err)
            }
            fmt.Println("__COMPLETE__[down]:" + "migration_file_name.go")
        }
    }
}
```

## Examples

See [examples.go](examples.go) for complete usage examples including:

- Basic usage
- Advanced configuration
- Custom filesystem usage
- Error handling patterns

## Compatibility

This library wraps the existing metana CLI functionality, so it maintains full compatibility with:

- All supported databases (PostgreSQL, MongoDB, file-based)
- Existing migration file formats
- All migration tracking features
- Dry-run capabilities

The CLI remains unchanged and continues to work as before.
