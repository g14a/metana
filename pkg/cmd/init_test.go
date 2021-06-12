package cmd

import (
	"bytes"
	"testing"

	"github.com/g14a/metana/pkg"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/spf13/cobra"
)

func Test_Init(t *testing.T) {
	var buf bytes.Buffer
	metanaCmd := NewMetanaCommand()
	initCmd := &cobra.Command{
		Use: "init",
		RunE: func(cmd *cobra.Command, args []string) error {
			FS := afero.NewMemMapFs()
			cmd.SetOut(&buf)
			return RunInit(cmd, args, FS, "/Users/g14a/metana")
		},
	}
	initCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	initCmd.Flags().StringP("env", "e", "", "Specify the environment to initialize migration")

	metanaCmd.AddCommand(initCmd)
	_, err := pkg.ExecuteCommand(metanaCmd, "init")
	assert.NoError(t, err)
	assert.Equal(t, " ✓ Created /Users/g14a/metana/migrations/main.go\n", buf.String())
}

func Test_Init_dir(t *testing.T) {
	var buf bytes.Buffer
	metanaCmd := NewMetanaCommand()
	initCmd := &cobra.Command{
		Use: "init",
		RunE: func(cmd *cobra.Command, args []string) error {
			FS := afero.NewMemMapFs()
			cmd.SetOut(&buf)
			return RunInit(cmd, args, FS, "/Users/g14a/metana")
		},
	}
	initCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	initCmd.Flags().StringP("env", "e", "", "Specify the environment to initialize migration")

	metanaCmd.AddCommand(initCmd)
	_, err := pkg.ExecuteCommand(metanaCmd, "init", "--dir=schema-mig")
	assert.NoError(t, err)
	assert.Equal(t, " ✓ Created /Users/g14a/metana/schema-mig/main.go\n", buf.String())
}

func Test_Init_config(t *testing.T) {
	var buf bytes.Buffer
	metanaCmd := NewMetanaCommand()
	initCmd := &cobra.Command{
		Use: "init",
		RunE: func(cmd *cobra.Command, args []string) error {
			FS := afero.NewMemMapFs()
			cmd.SetOut(&buf)
			afero.WriteFile(FS, "/Users/g14a/metana/.metana.yml", []byte("dir: schema-mig\nstore: \n"), 0644)
			return RunInit(cmd, args, FS, "/Users/g14a/metana")
		},
	}
	initCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	initCmd.Flags().StringP("env", "e", "", "Specify the environment to initialize migration")

	metanaCmd.AddCommand(initCmd)
	_, err := pkg.ExecuteCommand(metanaCmd, "init")
	assert.NoError(t, err)
	pkg.ExpectLines(t, buf.String(), []string{`✓ .metana.yml found`, ` ✓ Created \/Users\/g14a\/metana\/schema-mig\/*main.go`}...)
}

func NewMetanaCommand() *cobra.Command {
	metanaCmd := cobra.Command{
		Use: "metana",
	}
	return &metanaCmd
}
