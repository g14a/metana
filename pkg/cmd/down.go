package cmd

import (
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	migrate2 "github.com/g14a/metana/pkg/core/migrate"
	"github.com/g14a/metana/pkg/store"
	"github.com/g14a/metana/pkg/types"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunDown(cmd *cobra.Command, args []string, FS afero.Fs, wd string) {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		log.Fatal(err)
	}
	storeConn, err := cmd.Flags().GetString("store")
	if err != nil {
		log.Fatal(err)
	}

	dryRun, err := cmd.Flags().GetBool("dry")
	if err != nil {
		log.Fatal(err)
	}

	mc, _ := config.GetMetanaConfig(FS, wd)

	// Priority range is explicit, then config, then migrations
	var finalDir string

	if dir != "" {
		finalDir = dir
	} else if mc != nil && mc.Dir != "" && dir == "" {
		color.Green(" ✓ .metana.yml found")
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
		storeHouse, err = store.GetStoreViaConn(finalStoreConn, finalDir, FS)
		if err != nil {
			log.Fatal(err)
		}
		existingTrack, err = storeHouse.Load(FS)
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(existingTrack.Migrations) == 0 && !dryRun {
		color.Yellow("at least one upward migration needed")
		return
	} else {
		existingTrack.LastRunTS = int(time.Now().Unix())
	}

	upUntil, _ := cmd.Flags().GetString("until")
	errBuf := migrate2.Run(upUntil, finalDir, existingTrack.LastRunTS, false)

	if !dryRun {
		_, num := store.ProcessLogs(errBuf)
		track := store.TrackToSetDown(existingTrack, num)

		err = storeHouse.Set(track, FS)
		if err != nil {
			log.Fatal(err)
		}

		color.Green("  >>> migration : complete")
		return
	}
	color.White("  >>> dry run migration : complete")
}
