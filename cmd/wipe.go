// Package cmd /*
package cmd

import (
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	"github.com/g14a/metana/pkg/core/wipe"
	"github.com/spf13/cobra"
)

// wipeCmd represents the wipe command
var wipeCmd = &cobra.Command{
	Use:   "wipe",
	Short: "Wipe out old stale migration files and track in your store",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			log.Fatal(err)
		}

		store, err := cmd.Flags().GetString("store")
		if err != nil {
			log.Fatal(err)
		}

		mc, _ := config.GetMetanaConfig()

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
		if store != "" {
			finalStoreConn = store
		} else if mc != nil && mc.StoreConn != "" && store == "" {
			finalStoreConn = mc.StoreConn
		}

		confirmWipe := false

		prompt := &survey.Confirm{
			Message: "Wiping will delete stale migration files. Continue?",
		}
		survey.AskOne(prompt, &confirmWipe)

		if confirmWipe {
			// TODO
			err := wipe.Wipe(finalDir, finalStoreConn)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	wipeCmd.Flags().StringP("store", "s", "", "Specify a connection url to track migrations")
	wipeCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")

	rootCmd.AddCommand(wipeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// wipeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// wipeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
