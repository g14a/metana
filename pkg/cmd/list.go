package cmd

import (
	"log"
	"os"

	"github.com/g14a/metana/pkg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunList(cmd *cobra.Command, args []string, FS afero.Fs) {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		log.Fatal(err)
	}
	if dir == "" {
		dir = "migrations"
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = pkg.ListMigrations(wd, dir, FS)
	if err != nil {
		log.Fatal(err)
	}
}
