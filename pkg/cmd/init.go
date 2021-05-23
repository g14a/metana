package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	gen2 "github.com/g14a/metana/pkg/core/gen"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunInit(cmd *cobra.Command, args []string, FS afero.Fs, wd string) error {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		return err
	}

	mc, _ := config.GetMetanaConfig(FS, wd)
	// Priority range is explicit, then config, then migrations
	var finalDir string

	if dir != "" {
		finalDir = dir
		setMc := &config.MetanaConfig{
			Dir: finalDir,
		}
		err := config.SetMetanaConfig(setMc, FS, wd)
		if err != nil {
			return err
		}

	} else if mc != nil && mc.Dir != "" && dir == "" {
		fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ .metana.yml found\n"))
		finalDir = mc.Dir
		setMc := &config.MetanaConfig{
			Dir: finalDir,
		}
		err := config.SetMetanaConfig(setMc, FS, wd)
		if err != nil {
			return err
		}
	} else {
		finalDir = "migrations"
	}

	_ = FS.MkdirAll(finalDir+"/scripts", 0755)

	goModPath, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		return err
	}

	goModPathString := strings.TrimSpace(string(goModPath))
	if goModPathString == "" {
		color.Yellow("No go module found")
	}

	err = gen2.CreateInitConfig(finalDir, goModPathString, FS)
	if err != nil {
		return err
	}

	setMc := &config.MetanaConfig{
		Dir: finalDir,
	}

	if (&config.MetanaConfig{}) == mc || mc == nil {
		err := config.SetMetanaConfig(setMc, FS, wd)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(cmd.OutOrStdout(), color.GreenString(" ✓ Created "+wd+"/"+finalDir+"/main.go\n"))
	return nil
}
