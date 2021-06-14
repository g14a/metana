package cmd

import (
	"github.com/g14a/metana/pkg"
	"github.com/g14a/metana/pkg/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunList(cmd *cobra.Command, wd string, FS afero.Fs) error {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		return err
	}

	environment, err := cmd.Flags().GetString("env")
	if err != nil {
		return err
	}

	var finalDir string

	mc, _ := config.GetMetanaConfig(FS, wd)

	if dir != "" {
		finalDir = dir
	} else if mc != nil && mc.Dir != "" && dir == "" {
		finalDir = mc.Dir
	} else {
		finalDir = "migrations"
	}

	err = pkg.ListMigrations(cmd, wd, finalDir, FS, environment)
	if err != nil {
		return err
	}

	return nil
}
