// Package cmd /*
package cmd

import (
	"os"

	cmd2 "github.com/g14a/metana/pkg/cmd"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run upward migrations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		FS := afero.NewOsFs()
		wd, _ := os.Getwd()

		cmd2.RunUp(cmd, args, FS, wd)
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
