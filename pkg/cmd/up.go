package cmd

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	migrate2 "github.com/g14a/metana/pkg/core/migrate"
	"github.com/g14a/metana/pkg/store"
	"github.com/g14a/metana/pkg/types"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunUp(cmd *cobra.Command, args []string, FS afero.Fs, wd string) error {
	// Flags
	dir, _ := cmd.Flags().GetString("dir")
	storeConn, _ := cmd.Flags().GetString("store")
	until, _ := cmd.Flags().GetString("until")
	dryRun, _ := cmd.Flags().GetBool("dry")

	mc, _ := config.GetMetanaConfig(FS, wd)

	// Determine migration dir
	finalDir := "migrations"
	if dir != "" {
		finalDir = dir
	} else if mc != nil && mc.Dir != "" {
		fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ .metana.yml found\n"))
		finalDir = mc.Dir
	}

	// Determine store connection string
	finalStoreConn := ""
	if storeConn != "" {
		finalStoreConn = storeConn
	} else if mc != nil && mc.StoreConn != "" {
		finalStoreConn = mc.StoreConn
	}

	// Setup migration options
	opts := migrate2.MigrationOptions{
		Until:         until,
		MigrationsDir: finalDir,
		Wd:            wd,
		Up:            true,
		StoreConn:     finalStoreConn,
		DryRun:        dryRun,
	}

	// Always run migrations (dry or not)
	output, err := migrate2.Run(opts)
	track, _ := store.ProcessLogs(output)
	
	// Only persist to store if not a dry run
	if !dryRun && len(track.Migrations) > 0 {
		sh, err := store.GetStoreViaConn(finalStoreConn, finalDir, FS, wd)
		if err != nil {
			return err
		}

		existingTrack, err := sh.Load(FS)
		if err != nil {
			return err
		}

		existingTrack.LastRun = track.LastRun

		// Merge migrations (deduplicated by Title)
		existingMap := make(map[string]types.Migration)
		for _, m := range existingTrack.Migrations {
			existingMap[m.Title] = m
		}
		for _, m := range track.Migrations {
			existingMap[m.Title] = m
		}

		// Sort merged result by filename
		merged := make([]types.Migration, 0, len(existingMap))
		for _, m := range existingMap {
			merged = append(merged, m)
		}
		sort.Slice(merged, func(i, j int) bool {
			return merged[i].Title < merged[j].Title
		})

		existingTrack.Migrations = merged

		if err := sh.Set(existingTrack, FS); err != nil {
			return err
		}
	}

	if dryRun {
		fmt.Fprintf(cmd.OutOrStdout(), color.WhiteString("  >>> dry run migration : complete\n"))
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), color.GreenString("  >>> migration : complete\n"))
	}

	return err
}
