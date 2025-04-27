package cmd

import (
	"path/filepath"
	"testing"

	"github.com/g14a/metana/pkg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_Init_DefaultDir(t *testing.T) {
	t.Parallel()

	FS := afero.NewMemMapFs()
	tempDir := t.TempDir()

	cmd := NewMetanaCommand()
	initCmd := &cobra.Command{
		Use: "init",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunInit(cmd, FS, tempDir)
		},
	}
	initCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	cmd.AddCommand(initCmd)

	_, err := pkg.ExecuteCommand(cmd, "init")
	assert.NoError(t, err)

	// ✅ Check that migrations/scripts folder exists
	exists, err := afero.DirExists(FS, filepath.Join(tempDir, "migrations", "scripts"))
	assert.NoError(t, err)
	assert.True(t, exists)

	// ✅ Check .metana.yml is created
	exists, err = afero.Exists(FS, filepath.Join(tempDir, ".metana.yml"))
	assert.NoError(t, err)
	assert.True(t, exists)
}

func Test_Init_CustomDirFlag(t *testing.T) {
	t.Parallel()

	FS := afero.NewMemMapFs()
	tempDir := t.TempDir()

	cmd := NewMetanaCommand()
	initCmd := &cobra.Command{
		Use: "init",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunInit(cmd, FS, tempDir)
		},
	}
	initCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	cmd.AddCommand(initCmd)

	_, err := pkg.ExecuteCommand(cmd, "init", "--dir=schema-mig")
	assert.NoError(t, err)

	// ✅ Check that schema-mig/scripts folder exists
	exists, err := afero.DirExists(FS, filepath.Join(tempDir, "schema-mig", "scripts"))
	assert.NoError(t, err)
	assert.True(t, exists)

	// ✅ Check .metana.yml is created
	exists, err = afero.Exists(FS, filepath.Join(tempDir, ".metana.yml"))
	assert.NoError(t, err)
	assert.True(t, exists)
}

func Test_Init_WithExistingConfig(t *testing.T) {
	t.Parallel()

	FS := afero.NewMemMapFs()
	tempDir := t.TempDir()

	// Create .metana.yml manually
	err := afero.WriteFile(FS, filepath.Join(tempDir, ".metana.yml"), []byte("dir: schema-mig\nstore: \n"), 0644)
	assert.NoError(t, err)

	cmd := NewMetanaCommand()
	initCmd := &cobra.Command{
		Use: "init",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunInit(cmd, FS, tempDir)
		},
	}
	initCmd.Flags().StringP("dir", "d", "", "Specify custom migrations directory")
	cmd.AddCommand(initCmd)

	_, err = pkg.ExecuteCommand(cmd, "init")
	assert.NoError(t, err)

	// ✅ Check schema-mig/scripts folder created from config
	exists, err := afero.DirExists(FS, filepath.Join(tempDir, "schema-mig", "scripts"))
	assert.NoError(t, err)
	assert.True(t, exists)
}

func NewMetanaCommand() *cobra.Command {
	return &cobra.Command{
		Use: "metana",
	}
}
