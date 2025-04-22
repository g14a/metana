package cmd

import (
	"github.com/g14a/metana/pkg"
	"github.com/g14a/metana/pkg/config"
	"github.com/g14a/metana/pkg/store"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// RunList is the entrypoint for `metana list`
func RunList(cmd *cobra.Command, wd string, FS afero.Fs) error {
	dir, _ := cmd.Flags().GetString("dir")
	storeConn, _ := cmd.Flags().GetString("store")

	mc, _ := config.GetMetanaConfig(FS, wd)
	finalDir := resolveDir(dir, mc)
	finalStoreConn := resolveStore(storeConn, mc)

	// Always try to load the store — even if no explicit store connection string is provided
	st, err := store.GetStoreViaConn(finalStoreConn, finalDir, FS, wd)
	if err != nil {
		// fallback silently to no store — just show migrations without executed_at
		st = nil
	}

	return pkg.ListMigrations(cmd, wd, finalDir, FS, st)
}

func resolveDir(flagDir string, mc *config.MetanaConfig) string {
	if flagDir != "" {
		return flagDir
	}
	if mc != nil && mc.Dir != "" {
		return mc.Dir
	}
	return "migrations"
}

func resolveStore(flagStore string, mc *config.MetanaConfig) string {
	if flagStore != "" {
		return flagStore
	}
	if mc != nil && mc.StoreConn != "" {
		return mc.StoreConn
	}
	return "" // fallback to local migrate.json inside GetStoreViaConn
}
