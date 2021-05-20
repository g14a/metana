package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg"
	"github.com/g14a/metana/pkg/config"
	gen2 "github.com/g14a/metana/pkg/core/gen"
	"github.com/iancoleman/strcase"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunCreate(cmd *cobra.Command, args []string, FS afero.Fs, wd string) error {
	dir, err := cmd.Flags().GetString("dir")
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

	exists, err := gen2.MigrationExists(wd, finalDir, args[0], FS)
	if err != nil {
		return err
	}

	if exists {
		color.Yellow("Migration already exists")
		os.Exit(0)
	}

	firstMigration := false
	migrations, err := pkg.GetMigrations(wd, finalDir, FS)
	if err != nil {
		return err
	}

	if len(migrations) == 0 {
		firstMigration = true
	}

	fileName, err := gen2.CreateMigrationFile(finalDir, args[0], FS)
	if err != nil {
		color.Yellow("\nTry initializing migration using `metana init`\n\n")
		os.Exit(0)
	}

	err = gen2.Regen(finalDir, strcase.ToCamel(args[0]), strings.TrimPrefix(fileName, finalDir+"/scripts/"), firstMigration, FS)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ Created "+wd+"/"+fileName+"\n"))
	fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ Updated "+wd+"/"+finalDir+"/main.go\n"))

	return nil
}
