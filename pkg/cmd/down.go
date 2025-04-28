package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	"github.com/g14a/metana/pkg/core/migrate"
	"github.com/g14a/metana/pkg/store"
	"github.com/g14a/metana/pkg/types"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunDown(cmd *cobra.Command, args []string, fs afero.Fs, wd string) error {
	dirFlag, _ := cmd.Flags().GetString("dir")
	storeFlag, _ := cmd.Flags().GetString("store")
	until, _ := cmd.Flags().GetString("until")
	dryRun, _ := cmd.Flags().GetBool("dry")

	mc, _ := config.GetMetanaConfig(fs, wd)
	finalDir := resolveDir(dirFlag, mc)
	finalStore := resolveStore(storeFlag, mc)

	finalDir = cleanFinalDir(wd, finalDir)

	var track types.Track
	var storeHouse store.Store

	if !dryRun {
		sh, err := store.GetStoreViaConn(finalStore, finalDir, fs, wd)
		if err != nil {
			return err
		}
		track, err = sh.Load(fs)
		if err != nil {
			return err
		}
		storeHouse = sh
	}

	if len(track.Migrations) == 0 && !dryRun {
		fmt.Fprintf(cmd.OutOrStdout(), color.YellowString("at least one upward migration needed\n"))
		return nil
	}

	opts := migrate.MigrationOptions{
		Until:         until,
		MigrationsDir: finalDir,
		Wd:            wd,
		Up:            false,
		StoreConn:     finalStore,
		DryRun:        dryRun,
	}

	output, err := migrate.Run(opts)
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Fprintf(cmd.OutOrStdout(), color.WhiteString("  >>> dry run migration : complete\n"))
		return nil
	}

	_, num := store.ProcessLogs(output)
	newTrack := store.TrackToSetDown(track, num)

	if err := storeHouse.Set(newTrack, fs); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), color.GreenString("  >>> migration : complete\n"))
	return nil
}
