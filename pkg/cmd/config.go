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
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		return err
	}

	store, err := cmd.Flags().GetString("store")
	if err != nil {
		return err
	}

	env, err := cmd.Flags().GetString("env")
	if err != nil {
		return err
	}

	mc, err := config.GetMetanaConfig(FS, wd)
	if os.IsNotExist(err) {
		_, err = os.Create(".metana.yml")
		if err != nil {
			return err
		}
	}

	if env != "" {
		if len(mc.Environments) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), color.YellowString("No environment configured yet.\nTry initializing one with `metana init --env "+env+"`\n"))
			return nil
		}
		for i, e := range mc.Environments {
			if e.Name == env {
				if dir != "" {
					e.Dir = dir
				}
				if store != "" {
					e.Store = store
				}
				mc.Environments[i] = e
			}
		}
		if dir != "" {
			fmt.Fprintf(cmd.OutOrStdout(), color.YellowString(" ! Make sure you rename your exising environments directory to `"+dir+"`\n"))
		}
		err = config.SetEnvironmentMetanaConfig(mc, env, store, FS, wd)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ Set config\n"))
		return nil
	}

	if dir != "" {
		mc.Dir = dir
		fmt.Fprintf(cmd.OutOrStdout(), color.YellowString(" ! Make sure you rename your exising migrations directory to `"+dir+"`\n"))
	}

	if store != "" {
		mc.StoreConn = store
	}

	err = config.SetMetanaConfig(mc, FS, wd)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ Set config\n"))

	return nil
}
