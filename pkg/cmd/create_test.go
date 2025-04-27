package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/g14a/metana/pkg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_Create_Default(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	FS := afero.NewOsFs()

	// Step 1: Initialize migrations first
	initCmd := &cobra.Command{
		Use: "init",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunInit(cmd, FS, tempDir)
		},
	}
	initCmd.Flags().StringP("dir", "d", "", "")

	rootCmd := NewMetanaCommand()
	rootCmd.AddCommand(initCmd)

	_, err := pkg.ExecuteCommand(rootCmd, "init")
	assert.NoError(t, err)

	// Step 2: Now create a migration
	var buf bytes.Buffer
	createCmd := &cobra.Command{
		Use:  "create",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SetOut(&buf)
			return RunCreate(cmd, args, FS, tempDir)
		},
	}
	createCmd.Flags().StringP("dir", "d", "", "")
	createCmd.Flags().StringP("template", "t", "", "")

	rootCmd = NewMetanaCommand()
	rootCmd.AddCommand(createCmd)

	_, err = pkg.ExecuteCommand(rootCmd, "create", "abc")
	assert.NoError(t, err)

	// Step 3: Validate that migration file exists
	files, err := afero.ReadDir(FS, filepath.Join(tempDir, "migrations", "scripts"))
	assert.NoError(t, err)

	found := false
	for _, f := range files {
		if strings.Contains(strings.ToLower(f.Name()), "abc") && filepath.Ext(f.Name()) == ".go" {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func Test_Create_WithDirFlag(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	FS := afero.NewOsFs()

	// Step 1: Initialize migrations with custom dir
	initCmd := &cobra.Command{
		Use: "init",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set("dir", "schema-mig")
			return RunInit(cmd, FS, tempDir)
		},
	}
	initCmd.Flags().StringP("dir", "d", "", "")

	rootCmd := NewMetanaCommand()
	rootCmd.AddCommand(initCmd)

	_, err := pkg.ExecuteCommand(rootCmd, "init", "--dir=schema-mig")
	assert.NoError(t, err)

	// Step 2: Now create a migration
	var buf bytes.Buffer
	createCmd := &cobra.Command{
		Use:  "create",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SetOut(&buf)
			return RunCreate(cmd, args, FS, tempDir)
		},
	}
	createCmd.Flags().StringP("dir", "d", "", "")
	createCmd.Flags().StringP("template", "t", "", "")

	rootCmd = NewMetanaCommand()
	rootCmd.AddCommand(createCmd)

	_, err = pkg.ExecuteCommand(rootCmd, "create", "abc", "--dir=schema-mig")
	assert.NoError(t, err)

	// Step 3: Validate that migration file exists in schema-mig/scripts
	files, err := afero.ReadDir(FS, filepath.Join(tempDir, "schema-mig", "scripts"))
	assert.NoError(t, err)

	found := false
	for _, f := range files {
		if strings.Contains(strings.ToLower(f.Name()), "abc") && filepath.Ext(f.Name()) == ".go" {
			found = true
			break
		}
	}
	assert.True(t, found)
}
