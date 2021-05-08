// Package cmd /*
package cmd

import (
	"github.com/fatih/color"
	"github.com/g14a/go-migrate/pkg/migrate"
	"github.com/spf13/cobra"
	"log"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run upward migrations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			log.Fatal(err)
		}
		if dir == "" {
			dir = "migrations"
		}

		upUntil, _ := cmd.Flags().GetString("until")

		output, err := migrate.RunUp(upUntil, dir)

		if output != "" {
			color.Cyan("%v\n", output)
		}

		if err != nil {
			color.Red("  ERROR: %s", err)
		}

		color.Green("  >>> migration : complete")
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	upCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	upCmd.Flags().StringP("until", "u", "", "Migrate up until a specific point\n")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
