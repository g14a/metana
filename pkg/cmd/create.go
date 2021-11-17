package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/g14a/metana/pkg/core/environments"

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

	customTemplateFile, err := cmd.Flags().GetString("template")
	if err != nil {
		return err
	}

	environment, err := cmd.Flags().GetString("env")
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

	if environment != "" && !environments.CheckExistingEnvironment(FS, wd, finalDir, environment) {
		fmt.Fprintf(cmd.OutOrStdout(), color.YellowString("No environment configured yet.\nTry initializing one with `metana init --env "+environment+"`\n"))
		return nil
	}

	exists, err := gen2.MigrationExists(wd, finalDir, args[0], FS, environment)
	if err != nil {
		return err
	}

	if exists {
		color.Yellow("Migration already exists")
		os.Exit(0)
	}

	firstMigration := false
	migrations, err := pkg.GetMigrations(wd, finalDir, FS, environment)
	if err != nil {
		return err
	}

	if len(migrations) == 0 {
		firstMigration = true
	}

	opts := gen2.CreateMigrationOpts{
		Wd:            wd,
		MigrationsDir: finalDir,
		File:          args[0],
		CustomTmpl:    customTemplateFile,
		Environment:   environment,
		FS:            FS,
	}

	fileName, err := gen2.CreateMigrationFile(opts)
	if err != nil {
		color.Yellow("\nTry initializing migration using `metana init`\n\n")
		os.Exit(0)
	}

	var trimmedFile string

	if environment == "" {
		trimmedFile = strings.TrimPrefix(fileName, finalDir+"/scripts/")
	} else {
		trimmedFile = strings.TrimPrefix(fileName, finalDir+"/environments/"+environment+"/scripts/")
	}

	goModPath, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		return err
	}

	goModPathString := strings.TrimSpace(string(goModPath))
	if goModPathString == "" {
		color.Yellow("No go module found")
	}

	regenOpts := gen2.RegenOpts{
		MigrationsDir:  finalDir,
		MigrationName:  strcase.ToCamel(args[0]),
		Filename:       trimmedFile,
		FirstMigration: firstMigration,
		Environment:    environment,
		GoModPath:      goModPathString,
		Migrations:     migrations,
		FS:             FS,
	}

	err = gen2.Regen(regenOpts)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ Created "+wd+"/"+fileName+"\n"))
	if environment == "" {
		fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ Updated "+wd+"/"+finalDir+"/main.go\n"))
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ Updated "+wd+"/"+finalDir+"/environments/"+environment+"/main.go\n"))
	}

	return nil
}
