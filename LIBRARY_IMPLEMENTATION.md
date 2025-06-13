# Metana Library Interface - Implementation Summary

## Overview

I have successfully exposed the metana migration tool's functionality as a library that can be imported and used in other Go projects. This provides a programmatic interface to run migrations without requiring the CLI.

## What Was Created

### 1. Main Library Package (`pkg/metana/`)

- **`metana.go`** - Main library interface with the `Client` struct and all migration operations
- **`examples.go`** - Comprehensive usage examples
- **`metana_test.go`** - Test suite demonstrating library functionality
- **`README.md`** - Complete API documentation

### 2. Example Application (`example/`)

- **`main.go`** - Working example showing how to use the library
- **`README.md`** - Instructions for running the example

### 3. Updated Documentation

- Updated main `README.md` to include library usage section
- Added comprehensive documentation and examples

## Key Features

### ✅ **No Breaking Changes**
- Existing CLI functionality remains completely unchanged
- All existing commands work exactly as before
- Backward compatibility is maintained

### ✅ **Clean API Design**
- Struct-based configuration instead of CLI flags
- Consistent result handling with `MigrationResult`
- No Cobra/CLI dependencies in the library interface

### ✅ **Full Feature Parity**
- All CLI operations available: `Init`, `Create`, `Up`, `Down`, `List`
- Support for all store backends: PostgreSQL, MongoDB, file-based
- Dry-run capabilities
- Migration targeting with `Until` option
- Environment variable support for connection strings

### ✅ **Flexible Configuration**
- Custom working directories
- Custom migration directories
- Custom filesystem implementations (useful for testing)
- Optional store connection strings

### ✅ **Robust Error Handling**
- Structured error responses
- Detailed output capture
- Success/failure indicators

## Usage Examples

### Basic Usage

```go
import "github.com/g14a/metana/pkg/metana"

client := metana.NewClient(metana.Config{
    WorkingDirectory:      ".",
    MigrationsDirectory:   "migrations",
    StoreConnectionString: "postgres://user:pass@localhost/db",
})

// Initialize migration setup
result := client.Init()
if !result.Success {
    log.Fatal(result.Error)
}

// Create a migration
createResult := client.Create(metana.CreateOptions{
    Name: "add_users_table",
})

// Run migrations
upResult := client.Up(metana.UpOptions{})
```

### Advanced Usage

```go
// Run migrations with options
result := client.Up(metana.UpOptions{
    Until:  "specific_migration", // Stop at this migration
    DryRun: true,                // Preview without executing
    StoreConnectionString: "@DATABASE_URL", // Override default store
})

// Rollback migrations
downResult := client.Down(metana.DownOptions{
    Until: "initial_schema", // Rollback to this point
})

// List migrations
listResult := client.List(metana.ListOptions{})
```

### Testing with Custom Filesystem

```go
import "github.com/spf13/afero"

client := metana.NewClient(metana.Config{
    WorkingDirectory: "/tmp/test",
    FileSystem:       afero.NewMemMapFs(), // In-memory for testing
})
```

## Architecture

### How It Works

1. **Wrapper Pattern**: The library wraps the existing `pkg/cmd` functions
2. **Mock Cobra Commands**: Creates temporary Cobra commands to interface with existing code
3. **Output Capture**: Captures CLI output and converts it to structured results
4. **Error Translation**: Converts CLI errors to library-friendly error handling

### Code Structure

```
pkg/metana/
├── metana.go       # Main library interface
├── examples.go     # Usage examples
├── metana_test.go  # Test suite
└── README.md       # API documentation

example/
├── main.go         # Working example
└── README.md       # Example instructions
```

## API Reference

### Client Creation

```go
type Config struct {
    WorkingDirectory      string    // Base directory (default: current dir)
    MigrationsDirectory   string    // Migrations folder (default: "migrations")
    StoreConnectionString string    // DB connection (optional)
    FileSystem           afero.Fs  // Custom filesystem (optional)
}

client := metana.NewClient(config)
```

### Operations

All operations return `MigrationResult`:

```go
type MigrationResult struct {
    Success            bool              // Operation success
    Message            string            // Human-readable message
    Output             string            // Detailed output
    Error              error             // Any error that occurred
    ExecutedMigrations []types.Migration // Executed migrations list
}
```

Available operations:
- `client.Init()` - Initialize migration setup
- `client.Create(CreateOptions)` - Create new migration
- `client.Up(UpOptions)` - Run upward migrations
- `client.Down(DownOptions)` - Run downward migrations
- `client.List(ListOptions)` - List migrations
- `client.MigrationExists(name)` - Check if migration exists

## Testing

The library includes comprehensive tests:

```bash
cd /Users/gowtham.munukutla/metana
go test ./pkg/metana/...
```

## Running the Example

```bash
cd /Users/gowtham.munukutla/metana/example
go run main.go
```

This will demonstrate:
- Client initialization
- Migration setup
- Migration creation
- Migration listing
- Dry-run execution

## Integration Options

### 1. Application Startup Migrations

```go
func main() {
    client := metana.NewClient(metana.Config{
        StoreConnectionString: os.Getenv("DATABASE_URL"),
    })
    
    if result := client.Up(metana.UpOptions{}); !result.Success {
        log.Fatalf("Migration failed: %v", result.Error)
    }
    
    // Start your application
    startApplication()
}
```

### 2. Migration Management Service

```go
type MigrationService struct {
    client *metana.Client
}

func (s *MigrationService) RunMigrations() error {
    result := s.client.Up(metana.UpOptions{})
    return result.Error
}

func (s *MigrationService) GetMigrationStatus() ([]types.Migration, error) {
    result := s.client.List(metana.ListOptions{})
    return result.ExecutedMigrations, result.Error
}
```

### 3. Testing Integration

```go
func TestWithMigrations(t *testing.T) {
    client := metana.NewClient(metana.Config{
        FileSystem: afero.NewMemMapFs(),
        StoreConnectionString: "postgres://test:test@localhost/testdb",
    })
    
    // Setup test migrations
    client.Init()
    client.Up(metana.UpOptions{})
    
    // Run your tests
    
    // Cleanup
    client.Down(metana.DownOptions{})
}
```

## Benefits

1. **Programmatic Control**: Run migrations from your application code
2. **Integration Friendly**: Easy to integrate into existing applications
3. **Testing Support**: Use in-memory filesystems for testing
4. **No CLI Dependency**: Works without shell commands or CLI installation
5. **Type Safety**: Full Go type checking and IDE support
6. **Error Handling**: Proper Go error handling patterns
7. **Flexibility**: All CLI features available programmatically

## Compatibility

- ✅ Maintains full backward compatibility with CLI
- ✅ Supports all existing migration file formats
- ✅ Works with all supported databases (PostgreSQL, MongoDB, file-based)
- ✅ Preserves all existing features (dry-run, until, store options)
- ✅ Compatible with existing migration tracking

The library provides a clean, Go-idiomatic interface to metana's powerful migration capabilities while preserving all existing functionality.
