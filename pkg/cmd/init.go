package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunInit(cmd *cobra.Command, fs afero.Fs, wd string) error {
	dirFlag, _ := cmd.Flags().GetString("dir")
	mc, err := config.GetMetanaConfig(fs, wd)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if os.IsNotExist(err) {
		mc = &config.MetanaConfig{}
	}

	finalDir := resolveDir(dirFlag, mc)
	mc.Dir = finalDir

	if err := mkdirScripts(fs, finalDir); err != nil {
		return err
	}
	if err := config.SetMetanaConfig(mc, fs, finalDir); err != nil {
		return err
	}

	goModPath, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		return err
	}
	if strings.TrimSpace(string(goModPath)) == "" {
		color.Yellow("No go module found")
	}

	color.Green("Successfully initialized migration setup in %s", finalDir)
	return nil
}
