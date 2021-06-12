package cmd

import (
	"bytes"
	"testing"

	"github.com/g14a/metana/pkg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	var buf bytes.Buffer
	metanaCmd := NewMetanaCommand()
	createCmd := &cobra.Command{
		Use:  "create",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			FS := afero.NewMemMapFs()
			FS.MkdirAll("/Users/g14a/metana/migration/scripts", 0755)
			err := RunInit(cmd, args, FS, "/Users/g14a/metana")
			assert.NoError(t, err)
			cmd.SetOut(&buf)
			return RunCreate(cmd, args, FS, "/Users/g14a/metana")
		},
	}
	createCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	createCmd.Flags().StringP("template", "t", "", "Specify a custom Go template with Up and Down functions")
	createCmd.Flags().StringP("env", "e", "", "Specify an environment to create the migration")

	metanaCmd.AddCommand(createCmd)
	_, err := pkg.ExecuteCommand(metanaCmd, "create", "abc")
	assert.NoError(t, err)
	pkg.ExpectLines(t, buf.String(), []string{`✓ .metana.yml found`, `✓ Created \/Users\/g14a\/metana\/migrations\/scripts\/[0-9]*-Abc.go`, ` ✓ Updated \/Users\/g14a\/metana\/migrations\/*main.go`}...)
}

func Test_Create_dir(t *testing.T) {
	var buf bytes.Buffer
	metanaCmd := NewMetanaCommand()
	createCmd := &cobra.Command{
		Use:  "create",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			FS := afero.NewMemMapFs()
			FS.MkdirAll("/Users/g14a/metana/migration/scripts", 0755)
			err := RunInit(cmd, args, FS, "/Users/g14a/metana")
			assert.NoError(t, err)
			cmd.SetOut(&buf)
			return RunCreate(cmd, args, FS, "/Users/g14a/metana")
		},
	}
	createCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	createCmd.Flags().StringP("template", "t", "", "Specify a custom Go template with Up and Down functions")
	createCmd.Flags().StringP("env", "e", "", "Specify an environment to create the migration")

	metanaCmd.AddCommand(createCmd)
	_, err := pkg.ExecuteCommand(metanaCmd, "create", "abc", "--dir=schema-mig")
	assert.NoError(t, err)
	pkg.ExpectLines(t, buf.String(), []string{` ✓ Created \/Users\/g14a\/metana\/schema-mig\/scripts\/[0-9]*-Abc.go`, ` ✓ Updated \/Users\/g14a\/metana\/schema-mig\/*main.go`}...)
}

func Test_Create_Environment(t *testing.T) {
	var buf bytes.Buffer
	metanaCmd := NewMetanaCommand()
	createCmd := &cobra.Command{
		Use:  "create",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			FS := afero.NewMemMapFs()
			FS.MkdirAll("/Users/g14a/metana/schema-mig/scripts", 0755)
			err := RunInit(cmd, args, FS, "/Users/g14a/metana")
			assert.NoError(t, err)
			cmd.SetOut(&buf)
			return RunCreate(cmd, args, FS, "/Users/g14a/metana")
		},
	}
	createCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	createCmd.Flags().StringP("template", "t", "", "Specify a custom Go template with Up and Down functions")
	createCmd.Flags().StringP("env", "e", "", "Specify an environment to create the migration")

	metanaCmd.AddCommand(createCmd)
	_, err := pkg.ExecuteCommand(metanaCmd, "create", "abc", "--dir=schema-mig", "--env=dev")
	assert.NoError(t, err)
	pkg.ExpectLines(t, buf.String(), []string{` ✓ Created \/Users\/g14a\/metana\/schema-mig\/environments\/dev\/scripts\/[0-9]*-Abc.go`, ` ✓ Updated \/Users\/g14a\/metana\/schema-mig\/environments\/dev\/*main.go`}...)
}
