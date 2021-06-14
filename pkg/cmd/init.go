package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/g14a/metana/pkg/core/environments"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	gen2 "github.com/g14a/metana/pkg/core/gen"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunInit(cmd *cobra.Command, FS afero.Fs, wd string) error {
	dir, err := cmd.Flags().GetString("dir")
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

	setMc := &config.MetanaConfig{}

	switch {
	case environment != "":
		if dir == "" && mc != nil && mc.Dir != "" {
			dir = mc.Dir
			finalDir = dir
			setMc.Environments = mc.Environments
		} else if dir != "" {
			finalDir = dir
		} else {
			finalDir = "migrations"
		}
		envExists := environments.CheckExistingEnvironment(FS, wd, dir, environment)
		if envExists {
			fmt.Fprintf(cmd.OutOrStdout(), color.YellowString("Environment `"+environment+"` already exists\n"))
			return nil
		}
		environments.CheckExistingMigrationSetup(FS, wd)
		setMc.Dir = finalDir
		err := config.SetEnvironmentMetanaConfig(setMc, environment, "", FS, wd)
		if err != nil {
			return err
		}
		_ = FS.MkdirAll(finalDir+"/environments/"+environment+"/scripts", 0755)
	case dir != "":
		finalDir = dir
		setMc.Dir = finalDir
		err := config.SetMetanaConfig(setMc, FS, wd)
		if err != nil {
			return err
		}
		_ = FS.MkdirAll(finalDir+"/scripts", 0755)
	case mc != nil && mc.Dir != "" && dir == "":
		fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ .metana.yml found\n"))
		finalDir = mc.Dir
		setMc.Dir = finalDir
		setMc.Environments = mc.Environments
		err := config.SetMetanaConfig(setMc, FS, wd)
		if err != nil {
			return err
		}
		_ = FS.MkdirAll(finalDir+"/scripts", 0755)
	default:
		finalDir = "migrations"
		setMc.Dir = finalDir
		err := config.SetMetanaConfig(setMc, FS, wd)
		if err != nil {
			return err
		}
		_ = FS.MkdirAll(finalDir+"/scripts", 0755)
	}

	goModPath, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		return err
	}

	goModPathString := strings.TrimSpace(string(goModPath))
	if goModPathString == "" {
		color.Yellow("No go module found")
	}

	err = gen2.CreateInitConfig(finalDir, goModPathString, FS, environment)
	if err != nil {
		return err
	}

	if (&config.MetanaConfig{}) == mc || mc == nil {
		err := config.SetMetanaConfig(setMc, FS, wd)
		if err != nil {
			return err
		}
	}

	if environment == "" {
		fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ Created "+wd+"/"+finalDir+"/main.go\n"))
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ Created "+wd+"/"+finalDir+"/environments/"+environment+"/main.go\n"))
	}
	return nil
}
