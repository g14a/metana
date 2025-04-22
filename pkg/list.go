package pkg

import (
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/store"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func ListMigrations(cmd *cobra.Command, wd, migrationsDir string, fs afero.Fs, st store.Store) error {
	migrations, err := GetMigrations(wd, migrationsDir, fs)
	if err != nil {
		return err
	}

	executed := map[string]string{}
	if st != nil {
		track, err := st.Load(fs)
		if err != nil {
			color.Yellow("Warning: could not load migration store: %v", err)
		} else {
			for _, m := range track.Migrations {
				executed[m.Title] = m.ExecutedAt
			}
		}
	}

	var data [][]string
	if len(migrations) > 0 {
		for _, f := range migrations {
			execAt := ""
			if val, ok := executed[f.Name]; ok {
				execAt = val
			}
			data = append(data, []string{f.Name, execAt})
		}

		table := tablewriter.NewWriter(cmd.OutOrStdout())
		table.SetHeader([]string{"Migration", "Executed At"})

		for _, v := range data {
			table.Append(v)
		}
		table.Render()
	} else {
		color.Yellow("%s", "No migrations found")
	}

	return nil
}

func GetMigrations(wd, migrationsDir string, FS afero.Fs) ([]Migration, error) {
	matches, err := afero.Glob(FS, filepath.Join(wd, migrationsDir, "scripts", "[^.]*.*"))
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, f := range matches {
		if strings.Contains(f, "main.go") {
			continue
		}
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, Migration{
			Name: filepath.Base(f),
		})
	}

	return migrations, nil
}

type Migration struct {
	Name    string
	ModTime string
}
