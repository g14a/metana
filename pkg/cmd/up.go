package cmd

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	"github.com/g14a/metana/pkg/core/migrate"
	"github.com/g14a/metana/pkg/store"
	"github.com/g14a/metana/pkg/types"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunUp(cmd *cobra.Command, args []string, fs afero.Fs, wd string) error {
	dirFlag, _ := cmd.Flags().GetString("dir")
	storeFlag, _ := cmd.Flags().GetString("store")
	until, _ := cmd.Flags().GetString("until")
	dryRun, _ := cmd.Flags().GetBool("dry")

	mc, _ := config.GetMetanaConfig(fs, wd)
	finalDir := resolveDir(dirFlag, mc)
	finalStore := resolveStore(storeFlag, mc)

	finalDir = cleanFinalDir(wd, finalDir)

	opts := migrate.MigrationOptions{
		Until:         until,
		MigrationsDir: finalDir,
		Wd:            wd,
		Up:            true,
		StoreConn:     finalStore,
		DryRun:        dryRun,
	}

	output, err := migrate.Run(opts)
	track, _ := store.ProcessLogs(output)

	if !dryRun && len(track.Migrations) > 0 {
		sh, err := store.GetStoreViaConn(finalStore, finalDir, fs, wd)
		if err != nil {
			return err
		}

		existingTrack, err := sh.Load(fs)
		if err != nil {
			return err
		}

		existingMap := make(map[string]types.Migration)
		for _, m := range existingTrack.Migrations {
			existingMap[m.Title] = m
		}
		for _, m := range track.Migrations {
			existingMap[m.Title] = m
		}

		var merged []types.Migration
		for _, m := range existingMap {
			merged = append(merged, m)
		}
		sort.Slice(merged, func(i, j int) bool {
			return merged[i].Title < merged[j].Title
		})

		existingTrack.Migrations = merged
		existingTrack.LastRun = track.LastRun

		if err := sh.Set(existingTrack, fs); err != nil {
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
