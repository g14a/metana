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

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Run downward migrations",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		FS := afero.NewOsFs()
		wd, _ := os.Getwd()

		return cmd2.RunDown(cmd, args, FS, wd)
	},
}

func init() {
	rootCmd.AddCommand(downCmd)

	downCmd.Flags().StringP("until", "u", "", "Migrate down until a specific point")
	downCmd.Flags().StringP("store", "s", "", "Specify a connection URL to track migrations")
	downCmd.Flags().Bool("dry", false, "Specify if the downward migration is a dry run {true | false}")
	downCmd.Flags().StringP("env-file", "e", ".env", "Specify env file containing keys")

	downCmd.RegisterFlagCompletionFunc("until", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		FS := afero.NewOsFs()
		wd, err := os.Getwd()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		finalDir := "migrations"

		migrations, err := pkg.GetMigrations(filepath.Join(wd, finalDir), FS)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var names []string
		for _, m := range migrations {
			name := strings.TrimSuffix(m.Name, ".go")
			parts := strings.SplitN(name, "_", 2)
			if len(parts) == 2 {
				names = append(names, parts[1])
			}
		}

		return names, cobra.ShellCompDirectiveDefault
	})
}
