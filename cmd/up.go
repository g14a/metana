// Package cmd /*
package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/g14a/go-migrate/pkg/migrate"
	"github.com/g14a/go-migrate/pkg/store"
	"github.com/spf13/cobra"
	"log"
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

		if dir == "" {
			dir = "migrations"
		}

		storeHouse := store.GetStoreViaConn(storeConn, dir)
		existingTrack, err := storeHouse.Load()
		if err != nil {
			log.Fatal(err)
		}

		upUntil, _ := cmd.Flags().GetString("until")
		output, errBuf := migrate.RunUp(upUntil, dir, existingTrack.LastRunTS)

		if output != "" {
			color.Cyan("%v\n", output)
		}

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
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	upCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	upCmd.Flags().StringP("until", "u", "", "Migrate up until a specific point\n")
	upCmd.Flags().StringP("store", "s", "", "Specify a connection url to track migrations")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
