// Package cmd /*
package cmd

import (
	"log"
	"os"

	cmd2 "github.com/g14a/metana/pkg/cmd"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// InitCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a migrations directory",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		FS := afero.NewOsFs()
		wd, _ := os.Getwd()

		err := cmd2.RunInit(cmd, args, FS, wd)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	initCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	initCmd.Flags().StringP("env", "e", "", "Specify the environment to initialize migration")
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
