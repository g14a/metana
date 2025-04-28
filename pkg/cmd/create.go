package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	"github.com/g14a/metana/pkg/core/gen"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunCreate(cmd *cobra.Command, args []string, fs afero.Fs, wd string) error {
	time.Sleep(1 * time.Second)

	if len(args) == 0 {
		return fmt.Errorf("missing migration name")
	}

	dirFlag, _ := cmd.Flags().GetString("dir")
	mc, _ := config.GetMetanaConfig(fs, wd)
	finalDir := resolveDir(dirFlag, mc)

	if exists, _ := gen.MigrationExists(wd, finalDir, args[0], fs); exists {
		color.Yellow("Migration already exists")
		return nil
	}

	createdFile, err := gen.CreateMigrationFile(gen.CreateMigrationOpts{
		Wd:            wd,
		MigrationsDir: finalDir,
		File:          args[0],
		FS:            fs,
	})
	if err != nil {
		return fmt.Errorf("failed to create migration: %w\nTry initializing with `metana init`", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ Created %s\n", createdFile))
	return nil
}
