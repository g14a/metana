package cmd

import (
	"bytes"
	"testing"

	migrate2 "github.com/g14a/metana/pkg/core/migrate"

	"github.com/g14a/metana/pkg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_Down_AtleastOneMigrationNeeded(t *testing.T) {
	var buf bytes.Buffer
	metanaCmd := NewMetanaCommand()
	downCmd := &cobra.Command{
		Use: "down",
		RunE: func(cmd *cobra.Command, args []string) error {
			FS := afero.NewMemMapFs()
			FS.MkdirAll("/Users/g14a/metana/migration/scripts", 0755)
			err := RunInit(cmd, args, FS, "/Users/g14a/metana")
			assert.NoError(t, err)
			err = RunCreate(cmd, []string{"abc"}, FS, "/Users/g14a/metana")
			assert.NoError(t, err)
			cmd.SetOut(&buf)
			return RunDown(migrate2.MigrationOptions{
				MigrationsDir: "",
				Wd:            "/Users/g14a/metana",
				Up:            false,
				Cmd:           cmd,
			}, FS)
		},
	}
	downCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	downCmd.Flags().StringP("until", "u", "", "Migrate up until a specific point\n")
	downCmd.Flags().StringP("store", "s", "", "Specify a connection url to track migrations")
	downCmd.Flags().Bool("dry", false, "Specify if the upward migration is a dry run {true | false}")
	downCmd.Flags().StringP("template", "t", "", "Specify a custom Go template with Up and Down functions")
	downCmd.Flags().StringP("env", "e", ".env", "Specify environment keys from a file")

	metanaCmd.AddCommand(downCmd)
	_, err := pkg.ExecuteCommand(metanaCmd, "down")
	assert.NoError(t, err)
	pkg.ExpectLines(t, buf.String(), []string{`at least one upward migration needed`}...)

}
