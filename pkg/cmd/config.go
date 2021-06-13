package cmd

import (
	"log"
	"os"

	"github.com/g14a/metana/pkg/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunSetConfig(cmd *cobra.Command, FS afero.Fs, wd string) {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		log.Fatal(err)
	}

	store, err := cmd.Flags().GetString("store")
	if err != nil {
		log.Fatal(err)
	}

	mc, err := config.GetMetanaConfig(FS, wd)
	if os.IsNotExist(err) {
		_, err = os.Create(".metana.yml")
		if err != nil {
			log.Fatal(err)
		}
	}

	if dir != "" {
		mc.Dir = dir
	}

	if store != "" {
		mc.StoreConn = store
	}

	err = config.SetMetanaConfig(mc, FS, wd)
	if err != nil {
		log.Fatal(err)
	}
}
