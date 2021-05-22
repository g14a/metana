// Package cmd /*
package cmd

import (
	"log"
	"os"

	cmd2 "github.com/g14a/metana/pkg/cmd"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

// wipeCmd represents the wipe command
var wipeCmd = &cobra.Command{
	Use:   "wipe",
	Short: "Wipe out old stale migration files and track in your store",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		FS := afero.NewOsFs()
		wd, _ := os.Getwd()

		err := cmd2.RunWipe(cmd, args, FS, wd)
		if err != nil {
			log.Fatal(err)
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
