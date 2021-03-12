package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func ListMigrations() error {
	migrations, err := GetMigrations()
	if err != nil {
		return err
	}

	var data [][]string
	for _, f := range migrations {
		data = append(data, []string{f.Name, f.ModTime})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Migration", "Last Modified"})

	for _, v := range data {
		table.Append(v)
	}

	table.Render()

	return nil
}

func GetMigrations() ([]migration, error) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	m, err := filepath.Glob(wd + "/migrations/scripts/[^.]*.*")
	if err != nil {
		fmt.Println(err)
	}

	var migrations []migration
	for _, f := range m {
		fi, err := os.Stat(f)
		if err != nil {
			return nil, err
		}
		if strings.Contains(f, "main.go") || strings.Contains(f, "store.go") {
			continue
		}
		migrations = append(migrations, migration{
			Name:    filepath.Base(f),
			ModTime: fi.ModTime().Format("02-01-2006 15:04"),
		})
	}

	return migrations, nil
}

type migration struct {
	Name    string
	ModTime string
}
