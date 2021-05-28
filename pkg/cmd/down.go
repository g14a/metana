package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	migrate2 "github.com/g14a/metana/pkg/core/migrate"
	"github.com/g14a/metana/pkg/store"
	"github.com/g14a/metana/pkg/types"
	"github.com/spf13/afero"
)

func RunDown(opts migrate2.MigrationOptions, FS afero.Fs) error {
	var existingTrack types.Track
	var storeHouse store.Store
	if !opts.DryRun {
		sh, err := store.GetStoreViaConn(opts.StoreConn, opts.MigrationsDir, FS, opts.Wd)
		if err != nil {
			return err
		}
		existingTrack, err = sh.Load(FS)
		if err != nil {
			return err
		}
		storeHouse = sh
	}

	if len(existingTrack.Migrations) == 0 && !opts.DryRun {
		fmt.Fprintf(opts.Cmd.OutOrStdout(), color.YellowString("at least one upward migration needed\n"))
		return nil
	} else {
		existingTrack.LastRunTS = int(time.Now().Unix())
	}

	output, err := migrate2.Run(opts)
	if err != nil {
		return err
	}
	if !opts.DryRun {
		_, num := store.ProcessLogs(output)
		track := store.TrackToSetDown(existingTrack, num)

		err = storeHouse.Set(track, FS)
		if err != nil {
			return err
		}

		fmt.Fprintf(opts.Cmd.OutOrStdout(), color.GreenString("  >>> migration : complete\n"))
		return nil
	}
	fmt.Fprintf(opts.Cmd.OutOrStdout(), color.WhiteString("  >>> dry run migration : complete\n"))
	return nil
}
