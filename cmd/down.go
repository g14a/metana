// Package cmd /*
package cmd

import (
	"github.com/fatih/color"
	"github.com/g14a/go-migrate/pkg/migrate"
	"github.com/spf13/cobra"
	"log"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Run downward migrations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			log.Fatal(err)
		}
		if dir == "" {
			dir = "migrations"
		}

		downUntil, _ := cmd.Flags().GetString("until")
		output, err := migrate.RunDown(downUntil, dir)

		if err != nil {
			color.Red("\n  ERROR: %s", err)
		}

		if output != "" {
			color.Cyan("%v\n", output)
		}

		color.Green("  >>> migration : complete")
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
	downCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	downCmd.Flags().StringP("until", "u", "", "Migrate down until a specific point\n")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
