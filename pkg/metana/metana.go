// Package metana provides a programmatic interface to the metana migration tool.
// This package exposes the core migration functionality as a library that can be
// imported and used in other Go projects without requiring the CLI.
package metana

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/g14a/metana/pkg/cmd"
	"github.com/g14a/metana/pkg/core/gen"
	"github.com/g14a/metana/pkg/types"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// Client represents a metana migration client that can be used to perform
// migration operations programmatically.
type Client struct {
	config Config
	fs     afero.Fs
}

// Config holds the configuration for the metana migration client.
type Config struct {
	// WorkingDirectory is the base directory for migration operations.
	// If empty, uses the current working directory.
	WorkingDirectory string

	// MigrationsDirectory is the directory containing migration files.
	// Defaults to "migrations" if empty.
	MigrationsDirectory string

	// StoreConnectionString specifies the connection string for tracking migrations.
	// Supports postgres://, mongodb://, or file-based storage (default).
	// Can use environment variable references with @ENV_VAR_NAME syntax.
	StoreConnectionString string

	// FileSystem allows injecting a custom filesystem implementation.
	// Uses the OS filesystem if nil.
	FileSystem afero.Fs
}

// NewClient creates a new metana migration client with the provided configuration.
func NewClient(config Config) *Client {
	// Set defaults
	if config.WorkingDirectory == "" {
		if wd, err := os.Getwd(); err == nil {
			config.WorkingDirectory = wd
		} else {
			config.WorkingDirectory = "."
		}
	}

	if config.MigrationsDirectory == "" {
		config.MigrationsDirectory = "migrations"
	}

	fs := config.FileSystem
	if fs == nil {
		fs = afero.NewOsFs()
	}

	return &Client{
		config: config,
		fs:     fs,
	}
}

// UpOptions configures the behavior of upward migrations.
type UpOptions struct {
	// Until specifies a migration name to stop at.
	// If empty, runs all pending migrations.
	Until string

	// DryRun when true, shows what would be executed without actually running migrations.
	DryRun bool

	// StoreConnectionString overrides the client's default store connection.
	// If empty, uses the client's configured store connection.
	StoreConnectionString string
}

// DownOptions configures the behavior of downward migrations.
type DownOptions struct {
	// Until specifies a migration name to rollback to.
	// If empty, rolls back one migration.
	Until string

	// DryRun when true, shows what would be executed without actually running migrations.
	DryRun bool

	// StoreConnectionString overrides the client's default store connection.
	// If empty, uses the client's configured store connection.
	StoreConnectionString string
}

// CreateOptions configures the behavior of migration creation.
type CreateOptions struct {
	// Name is the name of the migration to create (required).
	Name string
}

// ListOptions configures the behavior of migration listing.
type ListOptions struct {
	// StoreConnectionString overrides the client's default store connection.
	// If empty, uses the client's configured store connection.
	StoreConnectionString string
}

// MigrationResult contains the result of a migration operation.
type MigrationResult struct {
	// Success indicates whether the operation completed successfully.
	Success bool

	// Message contains a human-readable result message.
	Message string

	// Output contains the detailed output from the migration execution.
	Output string

	// Error contains any error that occurred during the operation.
	Error error

	// ExecutedMigrations contains the list of migrations that were executed.
	ExecutedMigrations []types.Migration
}

// Up runs upward migrations according to the provided options.
func (c *Client) Up(opts UpOptions) MigrationResult {
	storeConn := opts.StoreConnectionString
	if storeConn == "" {
		storeConn = c.config.StoreConnectionString
	}

	// Create a mock cobra command to interface with existing pkg/cmd functions
	mockCmd := &cobra.Command{}
	mockCmd.SetOut(io.Discard)
	mockCmd.SetErr(io.Discard)

	// Set up flags that the existing RunUp function expects
	mockCmd.Flags().String("store", storeConn, "")
	mockCmd.Flags().String("until", opts.Until, "")
	mockCmd.Flags().Bool("dry", opts.DryRun, "")

	// Capture output
	var output bytes.Buffer
	mockCmd.SetOut(&output)

	// Call the existing RunUp function
	err := cmd.RunUp(mockCmd, []string{}, c.fs, c.config.WorkingDirectory)

	result := MigrationResult{
		Success: err == nil,
		Output:  output.String(),
		Error:   err,
	}

	if err == nil {
		if opts.DryRun {
			result.Message = "Dry run migration completed successfully"
		} else {
			result.Message = "Migration completed successfully"
		}
	} else {
		result.Message = fmt.Sprintf("Migration failed: %v", err)
	}

	return result
}

// Down runs downward migrations according to the provided options.
func (c *Client) Down(opts DownOptions) MigrationResult {
	storeConn := opts.StoreConnectionString
	if storeConn == "" {
		storeConn = c.config.StoreConnectionString
	}

	// Create a mock cobra command to interface with existing pkg/cmd functions
	mockCmd := &cobra.Command{}
	mockCmd.SetOut(io.Discard)
	mockCmd.SetErr(io.Discard)

	// Set up flags that the existing RunDown function expects
	mockCmd.Flags().String("store", storeConn, "")
	mockCmd.Flags().String("until", opts.Until, "")
	mockCmd.Flags().Bool("dry", opts.DryRun, "")

	// Capture output
	var output bytes.Buffer
	mockCmd.SetOut(&output)

	// Call the existing RunDown function
	err := cmd.RunDown(mockCmd, []string{}, c.fs, c.config.WorkingDirectory)

	result := MigrationResult{
		Success: err == nil,
		Output:  output.String(),
		Error:   err,
	}

	if err == nil {
		if opts.DryRun {
			result.Message = "Dry run rollback completed successfully"
		} else {
			result.Message = "Rollback completed successfully"
		}
	} else {
		result.Message = fmt.Sprintf("Rollback failed: %v", err)
	}

	return result
}

// Create creates a new migration file with the given name.
func (c *Client) Create(opts CreateOptions) MigrationResult {
	if opts.Name == "" {
		return MigrationResult{
			Success: false,
			Message: "Migration name is required",
			Error:   fmt.Errorf("migration name cannot be empty"),
		}
	}

	// Create a mock cobra command
	mockCmd := &cobra.Command{}
	var output bytes.Buffer
	mockCmd.SetOut(&output)

	// Call the existing RunCreate function
	err := cmd.RunCreate(mockCmd, []string{opts.Name}, c.fs, c.config.WorkingDirectory)

	result := MigrationResult{
		Success: err == nil,
		Output:  output.String(),
		Error:   err,
	}

	if err == nil {
		result.Message = fmt.Sprintf("Migration '%s' created successfully", opts.Name)
	} else {
		result.Message = fmt.Sprintf("Failed to create migration: %v", err)
	}

	return result
}

// Init initializes a new migration setup in the configured migrations directory.
func (c *Client) Init() MigrationResult {
	// Create a mock cobra command
	mockCmd := &cobra.Command{}
	var output bytes.Buffer
	mockCmd.SetOut(&output)

	// Call the existing RunInit function
	err := cmd.RunInit(mockCmd, c.fs, c.config.WorkingDirectory)

	result := MigrationResult{
		Success: err == nil,
		Output:  output.String(),
		Error:   err,
	}

	if err == nil {
		result.Message = fmt.Sprintf("Migration setup initialized in '%s'", c.config.MigrationsDirectory)
	} else {
		result.Message = fmt.Sprintf("Failed to initialize migration setup: %v", err)
	}

	return result
}

// List returns information about available and executed migrations.
func (c *Client) List(opts ListOptions) MigrationResult {
	storeConn := opts.StoreConnectionString
	if storeConn == "" {
		storeConn = c.config.StoreConnectionString
	}

	// Create a mock cobra command
	mockCmd := &cobra.Command{}
	mockCmd.Flags().String("store", storeConn, "")
	var output bytes.Buffer
	mockCmd.SetOut(&output)

	// Call the existing RunList function
	err := cmd.RunList(mockCmd, c.config.WorkingDirectory, c.fs)

	result := MigrationResult{
		Success: err == nil,
		Output:  output.String(),
		Error:   err,
	}

	if err == nil {
		result.Message = "Migration list retrieved successfully"
	} else {
		result.Message = fmt.Sprintf("Failed to list migrations: %v", err)
	}

	return result
}

// MigrationExists checks if a migration with the given name already exists.
func (c *Client) MigrationExists(name string) (bool, error) {
	migrationsDir := c.config.MigrationsDirectory
	if !filepath.IsAbs(migrationsDir) {
		migrationsDir = filepath.Join(c.config.WorkingDirectory, migrationsDir)
	}

	return gen.MigrationExists(c.config.WorkingDirectory, c.config.MigrationsDirectory, name, c.fs)
}

// GetWorkingDirectory returns the client's working directory.
func (c *Client) GetWorkingDirectory() string {
	return c.config.WorkingDirectory
}

// GetMigrationsDirectory returns the client's migrations directory.
func (c *Client) GetMigrationsDirectory() string {
	return c.config.MigrationsDirectory
}

// GetStoreConnectionString returns the client's store connection string.
func (c *Client) GetStoreConnectionString() string {
	return c.config.StoreConnectionString
}
