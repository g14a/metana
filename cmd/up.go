// Package cmd /*
package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/g14a/metana/pkg"

	cmd2 "github.com/g14a/metana/pkg/cmd"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run upward migrations",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		FS := afero.NewOsFs()
		wd, _ := os.Getwd()

		err := cmd2.RunUp(cmd, args, FS, wd)
		if err != nil {
			// Prevent Cobra from printing help on execution errors
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	upCmd.Flags().StringP("until", "u", "", "Migrate up until a specific point\n")
	upCmd.Flags().StringP("store", "s", "", "Specify a connection url to track migrations")
	upCmd.Flags().Bool("dry", false, "Specify if the upward migration is a dry run {true | false}")

	upCmd.RegisterFlagCompletionFunc("until", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		FS := afero.NewOsFs()

		wd, err := os.Getwd()
		if err != nil {
			return nil, 0
		}

		finalDir := "migrations"

		migrations, err := pkg.GetMigrations(filepath.Join(wd, finalDir), FS)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var names []string

		for _, m := range migrations {
			name := strings.TrimSuffix(m.Name, ".go")
			names = append(names, strings.Split(name, "-")[1])
		}

		return names, cobra.ShellCompDirectiveDefault
	})
}
