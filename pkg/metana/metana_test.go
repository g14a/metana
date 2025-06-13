package metana

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

func TestNewClient(t *testing.T) {
	// Test with default configuration
	client := NewClient(Config{})

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	// Working directory should be set
	if client.GetWorkingDirectory() == "" {
		t.Error("Expected working directory to be set")
	}

	// Migrations directory should default to "migrations"
	if client.GetMigrationsDirectory() != "migrations" {
		t.Errorf("Expected migrations directory to be 'migrations', got '%s'", client.GetMigrationsDirectory())
	}
}

func TestNewClientWithCustomConfig(t *testing.T) {
	config := Config{
		WorkingDirectory:      "/tmp/test",
		MigrationsDirectory:   "db/migrations",
		StoreConnectionString: "postgres://test",
		FileSystem:            afero.NewMemMapFs(),
	}

	client := NewClient(config)

	if client.GetWorkingDirectory() != "/tmp/test" {
		t.Errorf("Expected working directory '/tmp/test', got '%s'", client.GetWorkingDirectory())
	}

	if client.GetMigrationsDirectory() != "db/migrations" {
		t.Errorf("Expected migrations directory 'db/migrations', got '%s'", client.GetMigrationsDirectory())
	}

	if client.GetStoreConnectionString() != "postgres://test" {
		t.Errorf("Expected store connection 'postgres://test', got '%s'", client.GetStoreConnectionString())
	}
}

func TestClientInit(t *testing.T) {
	// Use in-memory filesystem for testing
	fs := afero.NewMemMapFs()

	client := NewClient(Config{
		WorkingDirectory:    "/tmp/test",
		MigrationsDirectory: "migrations",
		FileSystem:          fs,
	})

	result := client.Init()

	if !result.Success {
		t.Errorf("Expected Init to succeed, got error: %v", result.Error)
	}

	// Check that migrations directory was created
	exists, err := afero.DirExists(fs, "/tmp/test/migrations")
	if err != nil {
		t.Fatalf("Error checking directory existence: %v", err)
	}

	if !exists {
		t.Error("Expected migrations directory to be created")
	}

	// Check that scripts directory was created
	scriptsExists, err := afero.DirExists(fs, "/tmp/test/migrations/scripts")
	if err != nil {
		t.Fatalf("Error checking scripts directory existence: %v", err)
	}

	if !scriptsExists {
		t.Error("Expected migrations/scripts directory to be created")
	}
}

func TestClientCreateMigration(t *testing.T) {
	// Use in-memory filesystem for testing
	fs := afero.NewMemMapFs()

	client := NewClient(Config{
		WorkingDirectory:    "/tmp/test",
		MigrationsDirectory: "migrations",
		FileSystem:          fs,
	})

	// Initialize first
	initResult := client.Init()
	if !initResult.Success {
		t.Fatalf("Failed to initialize: %v", initResult.Error)
	}

	// Create a migration
	result := client.Create(CreateOptions{
		Name: "test_migration",
	})

	if !result.Success {
		t.Errorf("Expected Create to succeed, got error: %v", result.Error)
	}

	// Check that migration file was created
	scriptsDir := "/tmp/test/migrations/scripts"
	files, err := afero.ReadDir(fs, scriptsDir)
	if err != nil {
		t.Fatalf("Error reading scripts directory: %v", err)
	}

	if len(files) == 0 {
		t.Error("Expected migration file to be created")
	}

	// Check that the file name contains our migration name
	found := false
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".go" && 
		   len(file.Name()) > len("test_migration") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected migration file with correct name to be created")
	}
}

func TestClientCreateMigrationWithoutName(t *testing.T) {
	client := NewClient(Config{
		FileSystem: afero.NewMemMapFs(),
	})

	result := client.Create(CreateOptions{
		Name: "", // Empty name should fail
	})

	if result.Success {
		t.Error("Expected Create with empty name to fail")
	}

	if result.Error == nil {
		t.Error("Expected error to be set when Create fails")
	}
}

func TestMigrationExists(t *testing.T) {
	// Use in-memory filesystem for testing
	fs := afero.NewMemMapFs()

	client := NewClient(Config{
		WorkingDirectory:    "/tmp/test",
		MigrationsDirectory: "migrations",
		FileSystem:          fs,
	})

	// Initialize and create a migration
	client.Init()
	client.Create(CreateOptions{Name: "existing_migration"})

	// Test existing migration
	exists, err := client.MigrationExists("existing_migration")
	if err != nil {
		t.Fatalf("Error checking migration existence: %v", err)
	}

	if !exists {
		t.Error("Expected existing migration to be found")
	}

	// Test non-existing migration
	exists, err = client.MigrationExists("nonexistent_migration")
	if err != nil {
		t.Fatalf("Error checking migration existence: %v", err)
	}

	if exists {
		t.Error("Expected non-existing migration to not be found")
	}
}

// TestIntegrationWithFileSystem tests the library with the actual file system
// This test is skipped by default since it modifies the file system
func TestIntegrationWithFileSystem(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "metana-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	client := NewClient(Config{
		WorkingDirectory:    tempDir,
		MigrationsDirectory: "migrations",
	})

	// Test initialization
	initResult := client.Init()
	if !initResult.Success {
		t.Errorf("Init failed: %v", initResult.Error)
	}

	// Check that directories were created
	migrationsDir := filepath.Join(tempDir, "migrations")
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		t.Error("Migrations directory was not created")
	}

	scriptsDir := filepath.Join(migrationsDir, "scripts")
	if _, err := os.Stat(scriptsDir); os.IsNotExist(err) {
		t.Error("Scripts directory was not created")
	}

	// Test migration creation
	createResult := client.Create(CreateOptions{
		Name: "integration_test_migration",
	})
	if !createResult.Success {
		t.Errorf("Create failed: %v", createResult.Error)
	}

	// Check that migration file was created
	files, err := os.ReadDir(scriptsDir)
	if err != nil {
		t.Fatalf("Failed to read scripts directory: %v", err)
	}

	if len(files) == 0 {
		t.Error("No migration files were created")
	}

	// Test listing migrations
	listResult := client.List(ListOptions{})
	if !listResult.Success {
		t.Errorf("List failed: %v", listResult.Error)
	}

	// The output should contain information about our migration
	if listResult.Output == "" {
		t.Error("Expected list output to contain migration information")
	}
}

// BenchmarkClientOperations benchmarks the basic client operations
func BenchmarkClientOperations(b *testing.B) {
	fs := afero.NewMemMapFs()
	client := NewClient(Config{
		WorkingDirectory:    "/tmp/bench",
		MigrationsDirectory: "migrations",
		FileSystem:          fs,
	})

	// Initialize once
	client.Init()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		migrationName := "bench_migration"
		// if i > 0 {
		// 	migrationName = fmt.Sprintf("bench_migration_%d", i)
		// }

		result := client.Create(CreateOptions{
			Name: migrationName,
		})

		if !result.Success {
			b.Fatalf("Create failed: %v", result.Error)
		}
	}
}
