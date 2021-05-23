package cmd

import (
	"log"

	"github.com/g14a/metana/pkg"
	"github.com/g14a/metana/pkg/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunList(cmd *cobra.Command, args []string, wd string, FS afero.Fs) {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		log.Fatal(err)
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

	err = pkg.ListMigrations(wd, finalDir, FS)
	if err != nil {
		log.Fatal(err)
	}
}
