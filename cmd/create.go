// Package cmd /*
package cmd

import (
	"github.com/fatih/color"
	"github.com/g14a/go-migrate/pkg"
	"github.com/g14a/go-migrate/pkg/gen"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a migration in Go",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := cmd.Flags().GetString("dir")
		if len(args) == 0 || len(args) > 1 {
			color.Yellow("`create` accepts one argument")
			os.Exit(0)
		}

		if err != nil {
			log.Fatal(err)
		}
		if dir == "" {
			dir = "migrations"
		}

		exists, err := gen.MigrationExists(dir, args[0])
		if err != nil {
			log.Fatal(err)
		}

		if exists {
			color.Yellow("Migration already exists")
			os.Exit(0)
		}

		firstMigration := false
		migrations, err := pkg.GetMigrations(dir)
		if len(migrations) == 0 {
			firstMigration = true
		}

		fileName, err := gen.CreateMigrationFile(dir, args[0])
		if err != nil {
			color.Yellow("\nTry initializing migration using `go-migrate init`\n\n")
			os.Exit(0)
		}

		err = gen.Regen(dir, strcase.ToCamel(args[0]), strings.TrimPrefix(fileName, dir+"/scripts/"), firstMigration)
		if err != nil {
			log.Fatal(err)
		}

		wd, _ := os.Getwd()
		color.Green(" ✓ Created " + wd + "/" + fileName)
		color.Green(" ✓ Generated " + wd + "/" + dir + "/main.go")
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
