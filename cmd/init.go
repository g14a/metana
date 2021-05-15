// Package cmd /*
package cmd

import (
	"fmt"
	"log"
	"os"

	gen2 "github.com/g14a/metana/pkg/core/gen"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	"github.com/g14a/metana/pkg/initpkg"

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

		mc, _ := config.GetMetanaConfig()

		// Priority range is explicit, then config, then migrations
		var finalDir string

		fmt.Println("came here before dir==================")
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

		fmt.Println("came here after dir==========")

		goModPath, err := initpkg.GetGoModPath()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(goModPath,"==========gomodpath=========")

		err = gen2.CreateInitConfig(finalDir, goModPath)
		if err != nil {
			log.Fatal(err)
		}

		setMc := &config.MetanaConfig{
			Dir: finalDir,
		}

		if (&config.MetanaConfig{}) == mc || mc == nil {
			err := config.SetMetanaConfig(setMc)
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
