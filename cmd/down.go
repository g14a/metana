// Package cmd /*
package cmd

import (
	"os"

	cmd2 "github.com/g14a/metana/pkg/cmd"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Run downward migrations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		FS := afero.NewOsFs()
		wd, _ := os.Getwd()

		cmd2.RunDown(cmd, args, FS, wd)
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
	downCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	downCmd.Flags().StringP("until", "u", "", "Migrate down until a specific point\n")
	downCmd.Flags().StringP("store", "s", "", "Specify a connection url to track migrations")
	downCmd.Flags().Bool("dry", false, "Specify if the downward migration is a dry run {true | false}")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
