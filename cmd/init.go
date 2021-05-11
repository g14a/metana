// Package cmd /*
package cmd

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	"github.com/g14a/metana/pkg/gen"
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
		if mc.Dir != "" && dir == "" {
			dir = mc.Dir
		} else {
			dir = "migrations"
		}

		_ = os.MkdirAll(dir+"/scripts", 0755)
		wd, _ := os.Getwd()

		goModPath, err := initpkg.GetGoModPath()
		if err != nil {
			log.Fatal(err)
		}

		err = gen.CreateInitConfig(dir, goModPath)
		if err != nil {
			log.Fatal(err)
		}

		color.Green(" ✓ Created " + wd + "/" + dir + "/main.go")
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
