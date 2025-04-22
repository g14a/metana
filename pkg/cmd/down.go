package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	migrate2 "github.com/g14a/metana/pkg/core/migrate"
	"github.com/g14a/metana/pkg/store"
	"github.com/g14a/metana/pkg/types"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunDown(cmd *cobra.Command, args []string, FS afero.Fs, wd string) error {
	dir, _ := cmd.Flags().GetString("dir")
	storeConn, _ := cmd.Flags().GetString("store")
	until, _ := cmd.Flags().GetString("until")
	dryRun, _ := cmd.Flags().GetBool("dry")
	envFile, _ := cmd.Flags().GetString("env-file")

	mc, _ := config.GetMetanaConfig(FS, wd)

	// Resolve finalDir
	finalDir := "migrations"
	if dir != "" {
		finalDir = dir
	} else if mc != nil && mc.Dir != "" {
		fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ .metana.yml found\n"))
		finalDir = mc.Dir
	}

	// Resolve store connection
	finalStoreConn := ""
	if storeConn != "" {
		finalStoreConn = storeConn
	} else if mc != nil && mc.StoreConn != "" {
		finalStoreConn = mc.StoreConn
	}

	var existingTrack types.Track
	var storeHouse store.Store

	if !dryRun {
		sh, err := store.GetStoreViaConn(finalStoreConn, finalDir, FS, wd)
		if err != nil {
			return err
		}
		existingTrack, err = sh.Load(FS)
		if err != nil {
			return err
		}
		storeHouse = sh
	}

	if len(existingTrack.Migrations) == 0 && !dryRun {
		fmt.Fprintf(cmd.OutOrStdout(), color.YellowString("at least one upward migration needed\n"))
		return nil
	}

	opts := migrate2.MigrationOptions{
		Until:         until,
		MigrationsDir: finalDir,
		Wd:            wd,
		Up:            false,
		StoreConn:     finalStoreConn,
		DryRun:        dryRun,
		EnvFile:       envFile,
	}

	output, err := migrate2.Run(opts)
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Fprintf(cmd.OutOrStdout(), color.WhiteString("  >>> dry run migration : complete\n"))
		return nil
	}

	// Count how many __COMPLETE__ lines we got
	_, num := store.ProcessLogs(output)

	// Update track and remove `num` most recent migrations
	track := store.TrackToSetDown(existingTrack, num)
	err = storeHouse.Set(track, FS)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), color.GreenString("  >>> migration : complete\n"))
	return nil
}
