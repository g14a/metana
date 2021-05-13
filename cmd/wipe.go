// Package cmd /*
package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/g14a/metana/pkg/config"
	"github.com/spf13/cobra"
	"log"
)

// wipeCmd represents the wipe command
var wipeCmd = &cobra.Command{
	Use:   "wipe",
	Short: "Wipe out old stale migration files and track in your store",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		store, err := cmd.Flags().GetString("store")
		if err != nil {
			log.Fatal(err)
		}

		mc, _ := config.GetMetanaConfig()

		var finalStoreConn string
		if store != "" {
			finalStoreConn = store
		} else if mc != nil && mc.StoreConn != "" && store == "" {
			finalStoreConn = mc.StoreConn
		}

		fmt.Println(finalStoreConn)

		confirmWipe := false

		prompt := &survey.Confirm{
			Message: "Wiping will delete old migration files and existing store. Continue?",
		}
		survey.AskOne(prompt, &confirmWipe)

		if confirmWipe {
			// TODO
		}
	},
}

func init() {
	wipeCmd.Flags().StringP("store", "s", "", "Specify a connection url to track migrations")
	rootCmd.AddCommand(wipeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// wipeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// wipeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
