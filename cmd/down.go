// Package cmd /*
package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/migrate"
	"github.com/g14a/metana/pkg/store"
	"github.com/spf13/cobra"
	"log"
	"os"
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

		if dir == "" {
			dir = "migrations"
		}

		storeHouse := store.GetStoreViaConn(storeConn, dir)
		existingTrack, err := storeHouse.Load()
		if err != nil {
			log.Fatal(err)
		}

		if len(existingTrack.Migrations) == 0 {
			color.Yellow("at least one upward migration needed")
			os.Exit(0)
		}

		upUntil, _ := cmd.Flags().GetString("until")
		output, errBuf := migrate.RunDown(upUntil, dir, existingTrack.LastRunTS)

		if output != "" {
			color.Cyan("%v\n", output)
		}

		_, num := store.ProcessLogs(errBuf)
		track := store.TrackToSetDown(existingTrack, num)

		err = storeHouse.Set(track)
		if err != nil {
			fmt.Println(err)
		}

		color.Green("  >>> migration : complete")
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
	downCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	downCmd.Flags().StringP("until", "u", "", "Migrate down until a specific point\n")
	downCmd.Flags().StringP("store", "s", "", "Specify a connection url to track migrations")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
