// Package cmd /*
package cmd

import (
	"log"
	"os"

	cmd2 "github.com/g14a/metana/pkg/cmd"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a migration in Go",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		FS := afero.NewOsFs()
		wd, _ := os.Getwd()

		err := cmd2.RunCreate(cmd, args, FS, wd)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	createCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	createCmd.Flags().StringP("template", "t", "", "Specify a custom Go template with Up and Down functions")

	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
