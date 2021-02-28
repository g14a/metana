package pkg

import (
	"fmt"
	"os"
	"path/filepath"
)

func ListMigrations() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	m, err := filepath.Glob(wd + "/migrations/[^.]*.*")
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range m {
		fi, err := os.Stat(f)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(fi.ModTime())
	}
}