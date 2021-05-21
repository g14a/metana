package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/g14a/metana/pkg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_Up(t *testing.T) {
	var buf bytes.Buffer
	metanaCmd := NewMetanaCommand()
	upCmd := &cobra.Command{
		Use: "up",
		RunE: func(cmd *cobra.Command, args []string) error {
			FS := afero.NewMemMapFs()
			FS.MkdirAll("/Users/g14a/metana/migration/scripts", 0755)
			err := RunInit(cmd, args, FS, "/Users/g14a/metana")
			assert.NoError(t, err)
			err = RunCreate(cmd, []string{"abc"}, FS, "/Users/g14a/metana")
			assert.NoError(t, err)
			cmd.SetOut(&buf)
			return RunUp(cmd, []string{}, FS, "/Users/g14a/metana")
		},
	}
	upCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	upCmd.Flags().StringP("until", "u", "", "Migrate up until a specific point\n")
	upCmd.Flags().StringP("store", "s", "", "Specify a connection url to track migrations")
	upCmd.Flags().Bool("dry", false, "Specify if the upward migration is a dry run {true | false}")
	upCmd.Flags().StringP("template", "t", "", "Specify a custom Go template with Up and Down functions")

	metanaCmd.AddCommand(upCmd)
	_, err := pkg.ExecuteCommand(metanaCmd, "up")
	assert.NoError(t, err)
	pkg.ExpectLines(t, buf.String(), []string{`✓ .metana.yml found`, `  >>> migration : complete`}...)
}

func Test_Up_Dry(t *testing.T) {
	var buf bytes.Buffer
	metanaCmd := NewMetanaCommand()
	FS := afero.NewMemMapFs()

	upCmd := &cobra.Command{
		Use: "up",
		RunE: func(cmd *cobra.Command, args []string) error {
			FS.MkdirAll("/Users/g14a/metana/migration/scripts", 0755)
			err := RunInit(cmd, args, FS, "/Users/g14a/metana")
			assert.NoError(t, err)
			for _, m := range []string{"abc", "random", "addIndexes", "initSchema"} {
				err = RunCreate(cmd, []string{m}, FS, "/Users/g14a/metana")
				assert.NoError(t, err)
			}
			cmd.SetOut(&buf)
			return RunUp(cmd, []string{}, FS, "/Users/g14a/metana")
		},
	}
	upCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	upCmd.Flags().StringP("until", "u", "", "Migrate up until a specific point\n")
	upCmd.Flags().StringP("store", "s", "", "Specify a connection url to track migrations")
	upCmd.Flags().Bool("dry", false, "Specify if the upward migration is a dry run {true | false}")
	upCmd.Flags().StringP("template", "t", "", "Specify a custom Go template with Up and Down functions")
	metanaCmd.AddCommand(upCmd)
	_, err := pkg.ExecuteCommand(metanaCmd, "up", "--dry")
	assert.NoError(t, err)
	fmt.Println(buf.String())
	pkg.ExpectLines(t, buf.String(), []string{`✓ .metana.yml found`, `  >>> dry run migration : complete`}...)
}
