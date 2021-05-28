// Package cmd /*
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	migrate2 "github.com/g14a/metana/pkg/core/migrate"

	"github.com/g14a/metana/pkg"
	"github.com/g14a/metana/pkg/config"

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

		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}

		storeConn, err := cmd.Flags().GetString("store")
		if err != nil {
			return err
		}

		dryRun, err := cmd.Flags().GetBool("dry")
		if err != nil {
			return err
		}

		envFile, err := cmd.Flags().GetString("env")
		if err != nil {
			return err
		}

		mc, _ := config.GetMetanaConfig(FS, wd)

		// Priority range is explicit, then config, then migrations
		var finalDir string

		if dir != "" {
			finalDir = dir
		} else if mc != nil && mc.Dir != "" && dir == "" {
			fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ .metana.yml found\n"))
			finalDir = mc.Dir
		} else {
			finalDir = "migrations"
		}

		var finalStoreConn string
		if storeConn != "" {
			finalStoreConn = storeConn
		} else if mc != nil && mc.StoreConn != "" && storeConn == "" {
			finalStoreConn = mc.StoreConn
		}

		err = cmd2.RunDown(migrate2.MigrationOptions{
			MigrationsDir: finalDir,
			Wd:            wd,
			Up:            true,
			Cmd:           cmd,
			DryRun:        dryRun,
			StoreConn:     finalStoreConn,
			EnvFile:       envFile,
		}, FS)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
	downCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	downCmd.Flags().StringP("until", "u", "", "Migrate down until a specific point\n")
	downCmd.Flags().StringP("store", "s", "", "Specify a connection url to track migrations")
	downCmd.Flags().Bool("dry", false, "Specify if the downward migration is a dry run {true | false}")
	downCmd.Flags().StringP("env", "e", ".env", "Specify environment keys from a file")

	downCmd.RegisterFlagCompletionFunc("until", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		FS := afero.NewOsFs()

		wd, err := os.Getwd()
		if err != nil {
			return nil, 0
		}

		mc, _ := config.GetMetanaConfig(FS, wd)

		var finalDir string

		if mc != nil && mc.Dir != "" {
			finalDir = mc.Dir
		} else {
			finalDir = "migrations"
		}

		migrations, err := pkg.GetMigrations(wd, finalDir, FS)
		if err != nil {
			return nil, 0
		}

		var names []string

		for _, m := range migrations {
			name := strings.TrimSuffix(m.Name, ".go")
			names = append(names, strings.Split(name, "-")[1])
		}

		return names, cobra.ShellCompDirectiveDefault
	})
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
