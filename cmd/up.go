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

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run upward migrations",
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

		envFile, err := cmd.Flags().GetString("env-file")
		if err != nil {
			return err
		}

		env, err := cmd.Flags().GetString("env")
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

		err = cmd2.RunUp(migrate2.MigrationOptions{
			MigrationsDir: finalDir,
			Wd:            wd,
			Up:            true,
			Cmd:           cmd,
			DryRun:        dryRun,
			StoreConn:     finalStoreConn,
			Environment:   env,
			EnvFile:       envFile,
		}, FS)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	upCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	upCmd.Flags().StringP("until", "u", "", "Migrate up until a specific point\n")
	upCmd.Flags().StringP("store", "s", "", "Specify a connection url to track migrations")
	upCmd.Flags().Bool("dry", false, "Specify if the upward migration is a dry run {true | false}")
	upCmd.Flags().StringP("env-file", "e", ".env", "Specify file which contains env keys")
	upCmd.Flags().StringP("env", "", "", "Specify environment to run upward migration")

	upCmd.RegisterFlagCompletionFunc("until", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		environment, err := cmd.Flags().GetString("env")
		if err != nil {
			return nil, 0
		}

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

		migrations, err := pkg.GetMigrations(wd, finalDir, FS, environment)
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
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
