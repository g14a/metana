package cmd

import (
	"bytes"
	"testing"

	"github.com/g14a/metana/pkg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Set(t *testing.T) {

	var buf bytes.Buffer

	tests := []struct {
		args     []string
		function func() func(cmd *cobra.Command, args []string) error
		output   string
	}{
		{
			args: []string{"config", "set", "--store=random"},
			function: func() func(cmd *cobra.Command, args []string) error {
				return func(cmd *cobra.Command, args []string) error {
					FS := afero.NewMemMapFs()
					cmd.SetOut(&buf)
					afero.WriteFile(FS, "/Users/g14a/metana/.metana.yml", []byte("dir: schema-mig\nstore: \n"), 0644)
					err := RunSetConfig(cmd, FS, "/Users/g14a/metana")
					assert.NoError(t, err)
					file, err := afero.ReadFile(FS, "/Users/g14a/metana/.metana.yml")
					assert.NoError(t, err)
					assert.Equal(t, "dir: schema-mig\nstore: random\nenvironments: []\n", string(file))
					return nil
				}
			},
			output: " ✓ Set config\n",
		},
		{
			args: []string{"config", "set", "--dir=migrations"},
			function: func() func(cmd *cobra.Command, args []string) error {
				return func(cmd *cobra.Command, args []string) error {
					FS := afero.NewMemMapFs()
					cmd.SetOut(&buf)
					afero.WriteFile(FS, "/Users/g14a/metana/.metana.yml", []byte("dir: schema-mig\nstore: \n"), 0644)
					err := RunSetConfig(cmd, FS, "/Users/g14a/metana")
					assert.NoError(t, err)
					file, err := afero.ReadFile(FS, "/Users/g14a/metana/.metana.yml")
					assert.NoError(t, err)
					assert.Equal(t, "dir: migrations\nstore: \"\"\nenvironments: []\n", string(file))
					return nil
				}
			},
			output: " ! Make sure you rename your exising migrations directory to `migrations`\n ✓ Set config\n",
		},
		{
			args: []string{"config", "set", "--dir=migrations", "--env=dev"},
			function: func() func(cmd *cobra.Command, args []string) error {
				return func(cmd *cobra.Command, args []string) error {
					FS := afero.NewMemMapFs()
					cmd.SetOut(&buf)
					afero.WriteFile(FS, "/Users/g14a/metana/.metana.yml", []byte("dir: schema-mig\nstore: \n"), 0644)
					err := RunSetConfig(cmd, FS, "/Users/g14a/metana")
					assert.NoError(t, err)
					file, err := afero.ReadFile(FS, "/Users/g14a/metana/.metana.yml")
					assert.NoError(t, err)
					assert.Equal(t, "dir: schema-mig\nstore: \n", string(file))
					return nil
				}
			},
			output: "No environment configured yet.\nTry initializing one with `metana init --env dev`\n",
		},
		{
			args: []string{"config", "set", "--dir=change-env-dir", "--env=dev"},
			function: func() func(cmd *cobra.Command, args []string) error {
				return func(cmd *cobra.Command, args []string) error {
					FS := afero.NewMemMapFs()
					cmd.SetOut(&buf)
					afero.WriteFile(FS, "/Users/g14a/metana/.metana.yml", []byte("dir: schema-mig\nstore:\nenvironments:\n- name: dev\n  dir: dev\n  store: \"\"\n"), 0644)
					err := RunSetConfig(cmd, FS, "/Users/g14a/metana")
					assert.NoError(t, err)
					file, err := afero.ReadFile(FS, "/Users/g14a/metana/.metana.yml")
					assert.NoError(t, err)
					assert.Equal(t, "dir: schema-mig\nstore: \"\"\nenvironments:\n- name: dev\n  dir: change-env-dir\n  store: \"\"\n", string(file))
					return nil
				}
			},
			output: " ! Make sure you rename your exising environments directory to `change-env-dir`\n ✓ Set config\n",
		},
	}

	for _, tt := range tests {
		metanaCmd := NewMetanaCommand()

		configCmd := &cobra.Command{
			Use: "config",
			RunE: func(cmd *cobra.Command, args []string) error {
				return nil
			},
		}

		setCmd := &cobra.Command{
			Use:  "set",
			RunE: tt.function(),
		}
		setCmd.Flags().StringP("store", "s", "", "Set your store")
		setCmd.Flags().StringP("dir", "d", "", "Set your migrations directory")
		setCmd.Flags().StringP("env", "e", "", "Set config for your environment")
		configCmd.AddCommand(setCmd)
		metanaCmd.AddCommand(configCmd)
		c, out, err := pkg.ExecuteCommandC(metanaCmd, tt.args...)
		if out != "" {
			t.Errorf("Unexpected output: %v", out)
		}
		assert.NoError(t, err)
		assert.Equal(t, tt.output, buf.String())
		if c.Name() != "set" {
			t.Errorf(`invalid command returned from ExecuteC: expected "set"', got: %q`, c.Name())
		}
		buf.Reset()
	}
}
