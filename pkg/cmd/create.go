package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	gen2 "github.com/g14a/metana/pkg/core/gen"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunCreate(cmd *cobra.Command, args []string, FS afero.Fs, wd string) error {
	time.Sleep(1 * time.Second)

	if len(args) == 0 {
		return fmt.Errorf("missing migration name")
	}

	dir, _ := cmd.Flags().GetString("dir")
	templateFile, _ := cmd.Flags().GetString("template")

	mc, _ := config.GetMetanaConfig(FS, wd)
	finalDir := resolveMigrationsDir(dir, mc)

	if exists, _ := gen2.MigrationExists(wd, finalDir, args[0], FS); exists {
		color.Yellow("Migration already exists")
		return nil
	}

	// Create the migration file
	createOpts := gen2.CreateMigrationOpts{
		Wd:            wd,
		MigrationsDir: finalDir,
		File:          args[0],
		CustomTmpl:    templateFile,
		FS:            FS,
	}

	fileName, err := gen2.CreateMigrationFile(createOpts)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w\nTry initializing with `metana init`", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ Created %s/%s\n", wd, fileName))
	return nil
}

// Resolves the directory to use for migrations
func resolveMigrationsDir(dirFlag string, mc *config.MetanaConfig) string {
	if dirFlag != "" {
		return dirFlag
	}
	if mc != nil && mc.Dir != "" {
		color.Green(" ✓ .metana.yml found")
		return mc.Dir
	}
	return "migrations"
}
