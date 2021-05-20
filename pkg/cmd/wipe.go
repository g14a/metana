package cmd

import (
	"log"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/config"
	"github.com/g14a/metana/pkg/core/wipe"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunWipe(cmd *cobra.Command, args []string, FS afero.Fs, wd string) {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		log.Fatal(err)
	}

	store, err := cmd.Flags().GetString("store")
	if err != nil {
		log.Fatal(err)
	}

	mc, _ := config.GetMetanaConfig(FS, wd)

	var finalDir string

	if dir != "" {
		finalDir = dir
	} else if mc != nil && mc.Dir != "" && dir == "" {
		color.Green(" ✓ .metana.yml found")
		finalDir = mc.Dir
	} else {
		finalDir = "migrations"
	}

	var finalStoreConn string
	if store != "" {
		finalStoreConn = store
	} else if mc != nil && mc.StoreConn != "" && store == "" {
		finalStoreConn = mc.StoreConn
	}

	confirmWipe := false

	prompt := &survey.Confirm{
		Message: "Wiping will delete stale migration files. Continue?",
	}
	survey.AskOne(prompt, &confirmWipe)

	goModPath, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		log.Fatal(err)
	}

	goModPathString := strings.TrimSpace(string(goModPath))

	if confirmWipe {
		err := wipe.Wipe(goModPathString, wd, finalDir, finalStoreConn, FS)
		if err != nil {
			log.Fatal(err)
		}
	}
}
