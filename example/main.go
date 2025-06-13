// Package main demonstrates how to use metana as a library
package main

import (
	"fmt"
	"log"

	"github.com/g14a/metana/pkg/metana"
)

func main() {
	fmt.Println("=== Metana Library Example ===\n")

	// Create a metana client
	client := metana.NewClient(metana.Config{
		WorkingDirectory:      ".",
		MigrationsDirectory:   "example_migrations",
		StoreConnectionString: "", // Uses file-based storage
	})

	// 1. Initialize migration setup
	fmt.Println("1. Initializing migration setup...")
	initResult := client.Init()
	if !initResult.Success {
		log.Fatalf("Failed to initialize: %v", initResult.Error)
	}
	fmt.Printf("✓ %s\n\n", initResult.Message)

	// 2. Create a sample migration
	fmt.Println("2. Creating a sample migration...")
	createResult := client.Create(metana.CreateOptions{
		Name: "create_users_table",
	})
	if !createResult.Success {
		log.Fatalf("Failed to create migration: %v", createResult.Error)
	}
	fmt.Printf("✓ %s\n\n", createResult.Message)

	// 3. Check if migration exists
	fmt.Println("3. Checking if migration exists...")
	exists, err := client.MigrationExists("create_users_table")
	if err != nil {
		log.Printf("Error checking migration existence: %v", err)
	} else {
		fmt.Printf("✓ Migration exists: %v\n\n", exists)
	}

	// 4. List migrations
	fmt.Println("4. Listing migrations...")
	listResult := client.List(metana.ListOptions{})
	if !listResult.Success {
		log.Printf("Failed to list migrations: %v", listResult.Error)
	} else {
		fmt.Printf("✓ %s\n", listResult.Message)
		if listResult.Output != "" {
			fmt.Printf("Output:\n%s\n", listResult.Output)
		}
	}

	// 5. Demonstrate dry run
	fmt.Println("5. Running dry-run migration...")
	dryRunResult := client.Up(metana.UpOptions{
		DryRun: true,
	})
	if dryRunResult.Success {
		fmt.Printf("✓ %s\n", dryRunResult.Message)
	} else {
		fmt.Printf("⚠ Dry run failed (this is expected if no migration implementation exists): %s\n", dryRunResult.Message)
	}

	fmt.Println("\n=== Example completed ===")
	fmt.Println("Check the 'example_migrations' directory to see the created files.")
	fmt.Println("You can now implement your migration logic in the generated .go file.")
}
