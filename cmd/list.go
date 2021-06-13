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
	Run: func(cmd *cobra.Command, args []string) {
		FS := afero.NewOsFs()
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		cmd2.RunList(cmd, wd, FS)
	},
}

func init() {
	listCmd.Flags().StringP("dir", "d", "", "Specify migrations dir")
	listCmd.Flags().StringP("env", "e", "", "List migrations in an environment")

	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
