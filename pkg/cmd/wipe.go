package cmd

import (
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	"github.com/g14a/metana/pkg/core/wipe"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunWipe(cmd *cobra.Command, FS afero.Fs, wd string) error {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		return err
	}

	store, err := cmd.Flags().GetString("store")
	if err != nil {
		return err
	}

	confirmWipe, err := cmd.Flags().GetBool("yes")
	if err != nil {
		return err
	}

	environment, err := cmd.Flags().GetString("env")
	if err != nil {
		return err
	}

	mc, _ := config.GetMetanaConfig(FS, wd)

	var finalDir string

	if dir != "" {
		finalDir = dir
	} else if mc != nil && mc.Dir != "" && dir == "" {
		color.Green(" ✓ .metana.yml found")
		finalDir = mc.Dir
	} else {
		finalDir = "migrations"
	}

	var finalStoreConn string
	if store != "" {
		finalStoreConn = store
	} else if mc != nil && mc.StoreConn != "" && store == "" {
		finalStoreConn = mc.StoreConn
	}

	if !confirmWipe {
		prompt := &survey.Confirm{
			Message: "Wiping will delete stale migration files. Continue?",
		}
		survey.AskOne(prompt, &confirmWipe)
	}

	goModPath, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		return err
	}

	goModPathString := strings.TrimSpace(string(goModPath))

	if confirmWipe {
		opts := wipe.Opts{
			GoModPath:     goModPathString,
			Wd:            wd,
			MigrationsDir: finalDir,
			StoreConn:     finalStoreConn,
			Environment:   environment,
			FS:            FS,
		}
		err := wipe.Wipe(opts)
		if err != nil {
			return err
		}
	}

	return nil
}
