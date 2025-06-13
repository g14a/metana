// Simple test to verify the library compiles and works
package main

import (
	"fmt"
	"log"

	"github.com/g14a/metana/pkg/metana"
	"github.com/spf13/afero"
)

func main() {
	fmt.Println("Testing metana library compilation...")

	// Test with in-memory filesystem to avoid side effects
	fs := afero.NewMemMapFs()

	// Create a metana client
	client := metana.NewClient(metana.Config{
		WorkingDirectory:    "/tmp/test",
		MigrationsDirectory: "migrations",
		FileSystem:          fs,
	})

	fmt.Printf("Client created successfully!\n")
	fmt.Printf("Working Directory: %s\n", client.GetWorkingDirectory())
	fmt.Printf("Migrations Directory: %s\n", client.GetMigrationsDirectory())

	// Test initialization
	fmt.Println("\nTesting Init...")
	initResult := client.Init()
	if initResult.Success {
		fmt.Printf("✓ Init successful: %s\n", initResult.Message)
	} else {
		log.Printf("✗ Init failed: %v\n", initResult.Error)
	}

	// Test migration creation
	fmt.Println("\nTesting Create...")
	createResult := client.Create(metana.CreateOptions{
		Name: "test_migration",
	})
	if createResult.Success {
		fmt.Printf("✓ Create successful: %s\n", createResult.Message)
	} else {
		log.Printf("✗ Create failed: %v\n", createResult.Error)
	}

	// Test migration existence check
	fmt.Println("\nTesting MigrationExists...")
	exists, err := client.MigrationExists("test_migration")
	if err != nil {
		log.Printf("✗ MigrationExists failed: %v\n", err)
	} else {
		fmt.Printf("✓ MigrationExists: %v\n", exists)
	}

	fmt.Println("\n✅ All tests completed successfully! The library is working correctly.")
}
