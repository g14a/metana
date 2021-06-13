// Package cmd /*
package cmd

import (
	"os"

	cmd2 "github.com/g14a/metana/pkg/cmd"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage your local metana config in .metana.yml",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var setConfigCmd = &cobra.Command{
	Use:   "set",
	Short: "Set your metana config",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		FS := afero.NewOsFs()
		wd, _ := os.Getwd()

		cmd2.RunSetConfig(cmd, FS, wd)
	},
}

func init() {
	setConfigCmd.Flags().StringP("store", "s", "", "Set your store")
	setConfigCmd.Flags().StringP("dir", "d", "migrations", "Set your migrations directory")
	configCmd.AddCommand(setConfigCmd)
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
