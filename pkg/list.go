package pkg

import (
	"path/filepath"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/store"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func ListMigrations(cmd *cobra.Command, migrationsDir string, fs afero.Fs, st store.Store) error {
	migrations, err := GetMigrations(migrationsDir, fs)
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
	for _, f := range migrations {
		execAt := ""
		if val, ok := executed[f.Name]; ok {
			execAt = val
		}
		data = append(data, []string{f.Name, execAt})
	}

	if len(data) > 0 {
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

func GetMigrations(migrationsDir string, FS afero.Fs) ([]Migration, error) {

	matches, err := afero.Glob(FS, filepath.Join(migrationsDir, "scripts", "*.go"))
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, f := range matches {
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
