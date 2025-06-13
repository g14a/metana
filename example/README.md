# Metana Library Example

This example demonstrates how to use metana as a library in your Go applications.

## Running the Example

```bash
cd example
go run main.go
```

This will:

1. Initialize a migration setup in `example_migrations/`
2. Create a sample migration file
3. List the available migrations
4. Demonstrate a dry-run migration

## What You'll See

The example creates a directory structure like this:

```
example/
├── main.go
├── example_migrations/
│   ├── scripts/
│   │   └── [timestamp]_create_users_table.go
│   └── migrate.json
```

## Next Steps

After running the example:

1. Open the generated migration file in `example_migrations/scripts/`
2. Implement your `Up()` and `Down()` methods
3. Run actual migrations using the library

## Using in Your Project

Copy this pattern to use metana as a library in your own projects:

```go
import "github.com/g14a/metana/pkg/metana"

client := metana.NewClient(metana.Config{
    WorkingDirectory: ".",
    MigrationsDirectory: "migrations",
    StoreConnectionString: "postgres://user:pass@localhost/db", // optional
})

// Run migrations
result := client.Up(metana.UpOptions{})
if !result.Success {
    log.Fatal(result.Error)
}
```
