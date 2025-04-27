package cmd

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/g14a/metana/pkg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_Up(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		dryRun bool
		env    string
		output string
	}{
		{
			name:   "basic up",
			args:   []string{"up"},
			output: "  >>> migration : complete",
		},
		{
			name:   "dry run",
			args:   []string{"up", "--dry"},
			output: "  >>> dry run migration : complete",
			dryRun: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			// Setup memory FS and metana config
			FS := afero.NewMemMapFs()
			metanaPath := filepath.Join(tempDir, ".metana.yml")
			afero.WriteFile(FS, metanaPath, []byte(fmt.Sprintf("dir: %s\nstore: ''", tempDir)), 0644)

			cmd := &cobra.Command{
				Use: "up",
				RunE: func(cmd *cobra.Command, args []string) error {
					cmd.SetOut(&bytes.Buffer{})
					return RunUp(cmd, args, FS, tempDir)
				},
			}

			cmd.Flags().StringP("dir", "d", "", "")
			cmd.Flags().StringP("until", "u", "", "")
			cmd.Flags().StringP("store", "s", "", "")
			cmd.Flags().Bool("dry", false, "")
			cmd.Flags().StringP("env-file", "e", ".env", "")
			cmd.Flags().String("env", "", "")

			metanaCmd := NewMetanaCommand()
			metanaCmd.AddCommand(cmd)

			var buf bytes.Buffer
			cmd.SetOut(&buf)

			_, err := pkg.ExecuteCommand(metanaCmd, tt.args...)
			assert.NoError(t, err)
			assert.Contains(t, buf.String(), tt.output)
		})
	}
}
