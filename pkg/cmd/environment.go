package cmd

import (
	"github.com/fatih/color"
	environments2 "github.com/g14a/metana/pkg/core/environments"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunEnvironment(cmd *cobra.Command, args []string, FS afero.Fs, wd string) error {
	if environments2.CheckExistingMigrationSetup(FS, wd) {
		color.Yellow("Found existing migration setup.\nTry backing up your scripts and initialize a new migration setup.")
		return nil
	}
	return nil
}
