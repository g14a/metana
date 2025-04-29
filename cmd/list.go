// Package cmd /*
package cmd

import (
	"log"
	"os"

	cmd2 "github.com/g14a/metana/pkg/cmd"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List existing migrations",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		FS := afero.NewOsFs()
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		return cmd2.RunList(cmd, wd, FS)
	},
}

func init() {
	listCmd.Flags().StringP("store", "s", "", "Specify a connection URL to track migrations")

	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
