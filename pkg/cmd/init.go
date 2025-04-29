package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunInit(cmd *cobra.Command, fs afero.Fs, wd string) error {
	finalDir := resolveDir()

	if err := mkdirScripts(fs, finalDir); err != nil {
		return err
	}

	color.Green("Successfully initialized migration setup in '%s' folder", finalDir)
	return nil
}
