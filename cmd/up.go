// Package cmd /*
package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	"github.com/g14a/metana/pkg/migrate"
	"github.com/g14a/metana/pkg/store"
	"github.com/g14a/metana/pkg/types"
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run upward migrations",
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

		mc, _ := config.GetMetanaConfig()

		// Priority range is explicit, then config, then migrations
		if mc.Dir != "" && dir == "" {
			dir = mc.Dir
		} else {
			dir = "migrations"
		}

		if mc.StoreConn != "" && storeConn == "" {
			storeConn = mc.StoreConn
		}

		var existingTrack types.Track
		var storeHouse store.Store
		if !dryRun {
			storeHouse, err = store.GetStoreViaConn(storeConn, dir)
			if err != nil {
				log.Fatal(err)
			}
			existingTrack, err = storeHouse.Load()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			existingTrack.LastRunTS = 0
		}

		upUntil, _ := cmd.Flags().GetString("until")
		output, errBuf := migrate.RunUp(upUntil, dir, existingTrack.LastRunTS)

		if output != "" {
			color.Cyan("%v\n", output)
		}

		if !dryRun {
			track, _ := store.ProcessLogs(errBuf)

			existingTrack.LastRun = track.LastRun
			existingTrack.LastRunTS = track.LastRunTS
			existingTrack.Migrations = append(existingTrack.Migrations, track.Migrations...)

			if len(track.Migrations) > 0 {
				err = storeHouse.Set(existingTrack)
				if err != nil {
					fmt.Println(err)
				}
			}
			color.Green("  >>> migration : complete")
			return
		}
		color.White("  >>> dry run migration : complete")
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	upCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	upCmd.Flags().StringP("until", "u", "", "Migrate up until a specific point\n")
	upCmd.Flags().StringP("store", "s", "", "Specify a connection url to track migrations")
	upCmd.Flags().Bool("dry", false, "Specify if the upward migration is a dry run {true | false}")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
