package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunSetConfig(cmd *cobra.Command, FS afero.Fs, wd string) error {
	dir, _ := cmd.Flags().GetString("dir")
	store, _ := cmd.Flags().GetString("store")

	mc, err := config.GetMetanaConfig(FS, wd)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if os.IsNotExist(err) {
		mc = &config.MetanaConfig{}
	}

	// Global config update
	if dir != "" {
		mc.Dir = dir
		fmt.Fprintf(cmd.OutOrStdout(), color.YellowString(" ! Rename your migrations directory to `%s`\n", dir))
	}
	if store != "" {
		mc.StoreConn = store
	}

	if err := config.SetMetanaConfig(mc, FS, wd); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ Set config\n"))
	return nil
}
