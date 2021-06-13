package cmd

import (
	"fmt"
	"time"

	"github.com/g14a/metana/pkg/config"
	"github.com/spf13/cobra"

	"github.com/fatih/color"
	migrate2 "github.com/g14a/metana/pkg/core/migrate"
	"github.com/g14a/metana/pkg/store"
	"github.com/g14a/metana/pkg/types"
	"github.com/spf13/afero"
)

func RunDown(cmd *cobra.Command, args []string, FS afero.Fs, wd string) error {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		return err
	}

	storeConn, err := cmd.Flags().GetString("store")
	if err != nil {
		return err
	}

	until, err := cmd.Flags().GetString("until")
	if err != nil {
		return err
	}

	dryRun, err := cmd.Flags().GetBool("dry")
	if err != nil {
		return err
	}

	envFile, err := cmd.Flags().GetString("env-file")
	if err != nil {
		return err
	}

	env, err := cmd.Flags().GetString("env")
	if err != nil {
		return err
	}

	mc, _ := config.GetMetanaConfig(FS, wd)

	// Priority range is explicit, then config, then migrations
	var finalDir string

	if dir != "" {
		finalDir = dir
	} else if mc != nil && mc.Dir != "" && dir == "" {
		fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ .metana.yml found\n"))
		finalDir = mc.Dir
	} else {
		finalDir = "migrations"
	}

	var finalStoreConn string
	if storeConn != "" {
		finalStoreConn = storeConn
	} else if mc != nil && mc.StoreConn != "" && storeConn == "" {
		finalStoreConn = mc.StoreConn
	}

	var existingTrack types.Track
	var storeHouse store.Store
	if !dryRun {
		sh, err := store.GetStoreViaConn(finalStoreConn, finalDir, FS, wd, env)
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
	} else {
		existingTrack.LastRunTS = int(time.Now().Unix())
	}

	opts := migrate2.MigrationOptions{
		Until:         until,
		MigrationsDir: finalDir,
		Wd:            wd,
		LastRunTS:     existingTrack.LastRunTS,
		Up:            false,
		StoreConn:     finalStoreConn,
		DryRun:        dryRun,
		EnvFile:       envFile,
		Environment:   env,
	}

	output, err := migrate2.Run(opts)
	if err != nil {
		return err
	}
	if !dryRun {
		_, num := store.ProcessLogs(output)
		track := store.TrackToSetDown(existingTrack, num)

		err = storeHouse.Set(track, FS)
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), color.GreenString("  >>> migration : complete\n"))
		return nil
	}
	fmt.Fprintf(cmd.OutOrStdout(), color.WhiteString("  >>> dry run migration : complete\n"))
	return nil
}
