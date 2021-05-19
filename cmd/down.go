// Package cmd /*
package cmd

import (
	"log"
	"time"

	migrate2 "github.com/g14a/metana/pkg/core/migrate"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	"github.com/g14a/metana/pkg/store"
	"github.com/g14a/metana/pkg/types"
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Run downward migrations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
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

		mc, _ := config.GetMetanaConfig(FS)

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
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
	downCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	downCmd.Flags().StringP("until", "u", "", "Migrate down until a specific point\n")
	downCmd.Flags().StringP("store", "s", "", "Specify a connection url to track migrations")
	downCmd.Flags().Bool("dry", false, "Specify if the downward migration is a dry run {true | false}")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
