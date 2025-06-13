// Package metana provides examples for using the metana migration library.
package metana

import (
	"fmt"
	"log"
)

// ExampleBasicUsage demonstrates basic usage of the metana migration library.
func ExampleBasicUsage() {
	// Create a new metana client with default configuration
	client := NewClient(Config{
		WorkingDirectory:      "/path/to/your/project",
		MigrationsDirectory:   "migrations",                    // optional, defaults to "migrations"
		StoreConnectionString: "postgres://user:pass@host/db", // optional, defaults to file storage
	})

	// Initialize migration setup (creates migrations/scripts directory)
	initResult := client.Init()
	if !initResult.Success {
		log.Fatalf("Failed to initialize: %v", initResult.Error)
	}
	fmt.Println(initResult.Message)

	// Create a new migration
	createResult := client.Create(CreateOptions{
		Name: "add_users_table",
	})
	if !createResult.Success {
		log.Fatalf("Failed to create migration: %v", createResult.Error)
	}
	fmt.Println(createResult.Message)

	// Run upward migrations
	upResult := client.Up(UpOptions{
		DryRun: false, // set to true for dry run
	})
	if !upResult.Success {
		log.Fatalf("Migration failed: %v", upResult.Error)
	}
	fmt.Println(upResult.Message)

	// List migrations to see status
	listResult := client.List(ListOptions{})
	if listResult.Success {
		fmt.Println("Migration status:")
		fmt.Println(listResult.Output)
	}
}

// ExampleAdvancedUsage demonstrates advanced usage with custom options.
func ExampleAdvancedUsage() {
	// Create client with custom configuration
	client := NewClient(Config{
		WorkingDirectory:      ".",
		MigrationsDirectory:   "db/migrations",
		StoreConnectionString: "@DATABASE_URL", // Uses environment variable
	})

	// Run migrations up to a specific point
	result := client.Up(UpOptions{
		Until:  "add_indexes", // Stop at this migration
		DryRun: true,          // Preview what would happen
	})

	if result.Success {
		fmt.Printf("Dry run completed: %s\n", result.Message)
		fmt.Printf("Output:\n%s\n", result.Output)
	} else {
		fmt.Printf("Error: %v\n", result.Error)
	}

	// Rollback migrations
	downResult := client.Down(DownOptions{
		Until: "initial_schema", // Rollback to this migration
	})

	if downResult.Success {
		fmt.Printf("Rollback completed: %s\n", downResult.Message)
	}
}

// ExampleWithCustomFileSystem demonstrates using a custom filesystem.
func ExampleWithCustomFileSystem() {
	// You can use in-memory filesystem for testing or custom implementations
	// This example shows the interface, but you'd need to provide an actual
	// afero.Fs implementation

	client := NewClient(Config{
		WorkingDirectory:    "/tmp/test-project",
		MigrationsDirectory: "migrations",
		// FileSystem: afero.NewMemMapFs(), // Use in-memory filesystem for testing
	})

	// Check if a migration exists before creating it
	exists, err := client.MigrationExists("user_authentication")
	if err != nil {
		log.Printf("Error checking migration existence: %v", err)
		return
	}

	if !exists {
		result := client.Create(CreateOptions{
			Name: "user_authentication",
		})
		fmt.Printf("Create result: %s\n", result.Message)
	} else {
		fmt.Println("Migration already exists")
	}
}

// ExampleErrorHandling demonstrates proper error handling.
func ExampleErrorHandling() {
	client := NewClient(Config{
		WorkingDirectory: "/nonexistent/path",
	})

	// Always check the Success field and handle errors appropriately
	result := client.Up(UpOptions{})

	if !result.Success {
		// Log the error details
		log.Printf("Migration failed: %v", result.Error)
		log.Printf("Output: %s", result.Output)

		// Handle specific error types if needed
		switch {
		case result.Error != nil:
			// Handle the actual error
			fmt.Printf("Error occurred: %v\n", result.Error)
		default:
			// Handle other failure cases
			fmt.Printf("Migration failed: %s\n", result.Message)
		}
		return
	}

	fmt.Printf("Success: %s\n", result.Message)
}
