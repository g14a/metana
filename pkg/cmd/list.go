package cmd

import (
	"fmt"

	"github.com/g14a/metana/pkg"
	"github.com/g14a/metana/pkg/config"
	"github.com/g14a/metana/pkg/store"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunList(cmd *cobra.Command, wd string, fs afero.Fs) error {
	dirFlag, _ := cmd.Flags().GetString("dir")
	storeFlag, _ := cmd.Flags().GetString("store")

	mc, _ := config.GetMetanaConfig(fs, wd)
	finalDir := resolveDir(dirFlag, mc)
	finalStore := resolveStore(storeFlag, mc)

	st, err := store.GetStoreViaConn(finalStore, finalDir, fs, wd)
	if err != nil {
		fmt.Println("⚠️ Warning: store could not be initialized:", err)
		st = nil
	}

	return pkg.ListMigrations(cmd, finalDir, fs, st)
}
