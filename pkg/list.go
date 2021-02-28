package pkg

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"path/filepath"
	"strings"
)

func ListMigrations() error {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	m, err := filepath.Glob(wd + "/migrations/[^.]*.*")
	if err != nil {
		fmt.Println(err)
	}

	var data [][]string
	for _, f := range m {
		fi, err := os.Stat(f)
		if err != nil {
			return err
		}
		if strings.Contains(f, "main.go") {
			continue
		}

		data = append(data, []string{filepath.Base(f), fi.ModTime().Format("02-01-2006")})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Migration", "Last Modified"})

	for _, v := range data {
		table.Append(v)
	}

	table.Render()

	return nil
}
