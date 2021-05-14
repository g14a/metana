package pkg

import (
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func ListMigrations(migrationsDir string) error {
	migrations, err := GetMigrations(migrationsDir)
	if err != nil {
		return err
	}

	var data [][]string
	if len(migrations) > 0 {
		for _, f := range migrations {
			data = append(data, []string{f.Name, f.ModTime})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Migration", "Last Modified"})

		for _, v := range data {
			table.Append(v)
		}

		table.Render()
	} else {
		color.Yellow("%s", "No migrations found")
	}

	return nil
}

func GetMigrations(migrationsDir string) ([]Migration, error) {
	wd, err := os.Getwd()
	if err != nil {
		return []Migration{}, err
	}

	m, err := filepath.Glob(wd + "/" + migrationsDir + "/scripts/[^.]*.*")
	if err != nil {
		return []Migration{}, err
	}

	var migrations []Migration
	for _, f := range m {
		fi, err := os.Stat(f)
		if err != nil {
			return nil, err
		}
		if strings.Contains(f, "main.go") || strings.Contains(f, "store.go") {
			continue
		}
		migrations = append(migrations, Migration{
			Name:    filepath.Base(f),
			ModTime: fi.ModTime().Format("02-01-2006 15:04"),
		})
	}

	return migrations, nil
}

type Migration struct {
	Name    string
	ModTime string
}
