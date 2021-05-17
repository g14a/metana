// Package cmd /*
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/g14a/metana/pkg"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List existing migrations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			log.Fatal(err)
		}
		if dir == "" {
			dir = "migrations"
		}
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(wd, "===========outer wd==============")
		err = pkg.ListMigrations(wd, dir, FS)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	listCmd.Flags().StringP("dir", "d", "", "Specify migrations dir")
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
