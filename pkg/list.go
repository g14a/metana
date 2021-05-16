package pkg

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/fatih/color"

	"github.com/olekukonko/tablewriter"
)

func ListMigrations(migrationsDir string, fs afero.Fs) error {
	migrations, err := GetMigrations(migrationsDir, fs)
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

func GetMigrations(migrationsDir string, FS afero.Fs) ([]Migration, error) {
	FSUtil = &afero.Afero{Fs: FS}

	wd, err := os.Getwd()
	if err != nil {
		return []Migration{}, err
	}

	m, err := afero.Glob(FS, wd+"/"+migrationsDir+"/scripts/[^.]*.*")
	if err != nil {
		return []Migration{}, err
	}

	var migrations []Migration
	for _, f := range m {
		fi, err := FSUtil.Stat(f)
		if err != nil {
			return nil, err
		}
		if strings.Contains(f, "main.go") {
			continue
		}
		migrations = append(migrations, Migration{
			Name:    filepath.Base(f),
			ModTime: fi.ModTime().Format("02-01-2006 15:04"),
		})
	}

	return migrations, nil
}

func init() {
	FS = afero.NewOsFs()
	FSUtil = &afero.Afero{Fs: FS}
}

var (
	FS     afero.Fs
	FSUtil *afero.Afero
)

type Migration struct {
	Name    string
	ModTime string
}
