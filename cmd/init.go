// Package cmd /*
package cmd

import (
	"log"
	"os"
	"os/exec"
	"strings"

	gen2 "github.com/g14a/metana/pkg/core/gen"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a migrations directory",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			log.Fatal(err)
		}

		mc, _ := config.GetMetanaConfig(FS)

		// Priority range is explicit, then config, then migrations
		var finalDir string

		if dir != "" {
			finalDir = dir
		} else if mc != nil && mc.Dir != "" && dir == "" {
			color.Green(" ✓ .metana.yml found")
			finalDir = mc.Dir
		} else {
			finalDir = "migrations"
		}

		_ = os.MkdirAll(finalDir+"/scripts", 0755)
		wd, _ := os.Getwd()

		goModPath, err := exec.Command("go", "list", "-m").Output()
		if err != nil {
			log.Fatal(err)
		}

		goModPathString := strings.TrimSpace(string(goModPath))
		if goModPathString == "" {
			color.Yellow("No go module found")
		}

		err = gen2.CreateInitConfig(finalDir, goModPathString, FS)
		if err != nil {
			log.Fatal(err)
		}

		setMc := &config.MetanaConfig{
			Dir: finalDir,
		}

		if (&config.MetanaConfig{}) == mc || mc == nil {
			err := config.SetMetanaConfig(setMc, FS)
			if err != nil {
				return
			}
		}

		color.Green(" ✓ Created " + wd + "/" + finalDir + "/main.go")
	},
}

func init() {
	initCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
