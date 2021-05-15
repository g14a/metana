// Package cmd /*
package cmd

import (
	"log"
	"os"
	"strings"

	gen2 "github.com/g14a/metana/pkg/core/gen"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg"
	"github.com/g14a/metana/pkg/config"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a migration in Go",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			log.Fatal(err)
		}

		if len(args) == 0 || len(args) > 1 {
			color.Yellow("`create` accepts one argument")
			os.Exit(0)
		}

		mc, _ := config.GetMetanaConfig()

		// Priority range is explicit, then config, then migrations
		var finalDir string

		if dir != "" {
			finalDir = dir
		} else if mc != nil && mc.Dir != "" && dir == "" {
			color.Green(" ✓ .metana.yml found")
			finalDir = mc.Dir
		} else {
			finalDir = "migrations"
		}

		exists, err := gen2.MigrationExists(finalDir, args[0])
		if err != nil {
			log.Fatal(err)
		}

		if exists {
			color.Yellow("Migration already exists")
			os.Exit(0)
		}

		firstMigration := false
		migrations, err := pkg.GetMigrations(finalDir)
		if err != nil {
			log.Fatal(err)
		}

		if len(migrations) == 0 {
			firstMigration = true
		}

		fileName, err := gen2.CreateMigrationFile(finalDir, args[0])
		if err != nil {
			color.Yellow("\nTry initializing migration using `metana init`\n\n")
			os.Exit(0)
		}

		err = gen2.Regen(finalDir, strcase.ToCamel(args[0]), strings.TrimPrefix(fileName, finalDir+"/scripts/"), firstMigration)
		if err != nil {
			log.Fatal(err)
		}

		wd, _ := os.Getwd()
		color.Green(" ✓ Created " + wd + "/" + fileName)
		color.Green(" ✓ Updated " + wd + "/" + finalDir + "/main.go")
	},
}

func init() {
	createCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	rootCmd.AddCommand(createCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
