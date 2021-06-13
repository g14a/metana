package cmd

import (
	"bytes"
	"testing"

	"github.com/g14a/metana/pkg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_Up(t *testing.T) {
	var buf bytes.Buffer

	tests := []struct {
		cmd    *cobra.Command
		args   []string
		output []string
	}{
		{
			cmd: &cobra.Command{
				Use: "up",
				RunE: func(cmd *cobra.Command, args []string) error {
					FS := afero.NewOsFs()
					cmd.SetOut(&buf)
					afero.WriteFile(FS, "../../.metana.yml", []byte("dir: testdata\nstore: ''"), 0644)
					return RunUp(cmd, args, FS, "../..")
				},
			},
			args:   []string{"up"},
			output: []string{`  >>> migration : complete`},
		},
		{
			cmd: &cobra.Command{
				Use: "up",
				RunE: func(cmd *cobra.Command, args []string) error {
					FS := afero.NewOsFs()
					cmd.SetOut(&buf)
					afero.WriteFile(FS, "../../.metana.yml", []byte("dir: testdata\nstore: ''"), 0644)
					return RunUp(cmd, args, FS, "../..")
				},
			},
			args:   []string{"up", "--dry"},
			output: []string{`  >>> dry run migration : complete`},
		},
		{
			cmd: &cobra.Command{
				Use: "up",
				RunE: func(cmd *cobra.Command, args []string) error {
					FS := afero.NewOsFs()
					cmd.SetOut(&buf)
					afero.WriteFile(FS, "../../.metana.yml", []byte("dir: testdata\nstore: ''"), 0644)
					return RunUp(cmd, args, FS, "../..")
				},
			},
			args:   []string{"up", "--env=dev"},
			output: []string{`  >>> migration : complete`},
		},
	}

	for _, tt := range tests {
		tt.cmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
		tt.cmd.Flags().StringP("until", "u", "", "Migrate up until a specific point\n")
		tt.cmd.Flags().StringP("store", "s", "", "Specify a connection url to track migrations")
		tt.cmd.Flags().Bool("dry", false, "Specify if the upward migration is a dry run {true | false}")
		tt.cmd.Flags().StringP("template", "t", "", "Specify a custom Go template with Up and Down functions")
		tt.cmd.Flags().StringP("env", "", "", "Specify environment keys from a file")
		tt.cmd.Flags().StringP("env-file", "e", ".env", "Specify file which contains env keys")
		metanaCmd := NewMetanaCommand()
		metanaCmd.AddCommand(tt.cmd)
		_, err := pkg.ExecuteCommand(metanaCmd, tt.args...)
		assert.NoError(t, err)
		pkg.ExpectLines(t, buf.String(), tt.output...)
	}
}
