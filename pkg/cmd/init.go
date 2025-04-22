package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunInit(cmd *cobra.Command, FS afero.Fs, wd string) error {
	dir, _ := cmd.Flags().GetString("dir")

	// Load existing config or create new
	mc, err := config.GetMetanaConfig(FS, wd)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if os.IsNotExist(err) {
		mc = &config.MetanaConfig{}
	}

	// Resolve finalDir: CLI > config > default
	finalDir := "migrations"
	if dir != "" {
		finalDir = dir
	} else if mc.Dir != "" {
		finalDir = mc.Dir
	}
	mc.Dir = finalDir

	// Create scripts directory
	if err := FS.MkdirAll(fmt.Sprintf("%s/scripts", finalDir), 0755); err != nil {
		return err
	}

	// Save updated config
	if err := config.SetMetanaConfig(mc, FS, wd); err != nil {
		return err
	}

	// Determine Go module path
	goModPath, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		return err
	}
	goMod := strings.TrimSpace(string(goModPath))
	if goMod == "" {
		color.Yellow("No go module found")
	}
	color.Green("Successfully initialized migration setup in %s", finalDir)
	return nil
}
