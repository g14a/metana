package pkg

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/fatih/color"

	"github.com/olekukonko/tablewriter"
)

func ListMigrations(wd, migrationsDir string, fs afero.Fs, environment string) error {
	migrations, err := GetMigrations(wd, migrationsDir, fs, environment)
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

func GetMigrations(wd, migrationsDir string, FS afero.Fs, environment string) ([]Migration, error) {
	FSUtil := &afero.Afero{Fs: FS}
	var m []string
	var err error
	if environment == "" {
		m, err = afero.Glob(FS, wd+"/"+migrationsDir+"/scripts/[^.]*.*")
		if err != nil {
			return []Migration{}, err
		}
	} else {
		m, err = afero.Glob(FS, wd+"/"+migrationsDir+"/environments/"+environment+"/scripts/[^.]*.*")
		if err != nil {
			return []Migration{}, err
		}
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

type Migration struct {
	Name    string
	ModTime string
}
