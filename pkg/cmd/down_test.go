package cmd

import (
	"bytes"
	"testing"

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
			FS.MkdirAll("/Users/g14a/metana/migrations/scripts", 0755)
			err := RunInit(cmd, FS, "/Users/g14a/metana")
			assert.NoError(t, err)
			err = RunCreate(cmd, []string{"abc"}, FS, "/Users/g14a/metana")
			assert.NoError(t, err)
			cmd.SetOut(&buf)
			return RunDown(cmd, args, FS, "migrations")
		},
	}
	downCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	downCmd.Flags().StringP("until", "u", "", "Migrate up until a specific point\n")
	downCmd.Flags().StringP("store", "s", "", "Specify a connection url to track migrations")
	downCmd.Flags().Bool("dry", false, "Specify if the upward migration is a dry run {true | false}")
	downCmd.Flags().StringP("template", "t", "", "Specify a custom Go template with Up and Down functions")
	downCmd.Flags().StringP("env-file", "e", ".env", "Specify file which contains env keys")
	downCmd.Flags().StringP("env", "", "", "Specify environment to run downward migration")

	metanaCmd.AddCommand(downCmd)
	_, err := pkg.ExecuteCommand(metanaCmd, "down")
	assert.NoError(t, err)
	pkg.ExpectLines(t, buf.String(), []string{`at least one upward migration needed`}...)

}
