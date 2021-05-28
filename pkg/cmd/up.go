package cmd

import (
	"fmt"

	"github.com/fatih/color"
	migrate2 "github.com/g14a/metana/pkg/core/migrate"
	"github.com/g14a/metana/pkg/store"
	"github.com/g14a/metana/pkg/types"
	"github.com/spf13/afero"
)

func RunUp(opts migrate2.MigrationOptions, FS afero.Fs, wd string) error {

	var existingTrack types.Track
	var storeHouse store.Store
	if !opts.DryRun {
		storeHouse, err := store.GetStoreViaConn(opts.StoreConn, opts.MigrationsDir, FS, wd)
		if err != nil {
			return err
		}
		existingTrack, err = storeHouse.Load(FS)
		if err != nil {
			return err
		}
	} else {
		existingTrack.LastRunTS = 0
	}

	output, err := migrate2.Run(opts.Until, opts.MigrationsDir, wd, existingTrack.LastRunTS, true)
	if err != nil {
		return err
	}

	if !opts.DryRun {
		track, _ := store.ProcessLogs(output)

		existingTrack.LastRun = track.LastRun
		existingTrack.LastRunTS = track.LastRunTS
		existingTrack.Migrations = append(existingTrack.Migrations, track.Migrations...)

		if len(track.Migrations) > 0 {
			err = storeHouse.Set(existingTrack, FS)
			if err != nil {
				return err
			}
		}
		fmt.Fprintf(opts.Cmd.OutOrStdout(), color.GreenString("  >>> migration : complete\n"))

		return nil
	}
	fmt.Fprintf(opts.Cmd.OutOrStdout(), color.WhiteString("  >>> dry run migration : complete\n"))
	return nil
}
