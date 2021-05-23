package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/mitchellh/go-homedir"

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
			FS := afero.NewOsFs()
			_, err := homedir.Dir()
			assert.NoError(t, err)
			cmd.SetOut(&buf)
			wd, _ := os.Getwd()
			fmt.Println(wd, "===========wd===========")
			afero.WriteFile(FS, "../../.metana.yml", []byte("dir: testdata\nstore: ''"), 0644)
			return RunUp(cmd, []string{}, FS, "../..")
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

	upCmd := &cobra.Command{
		Use: "up",
		RunE: func(cmd *cobra.Command, args []string) error {
			FS := afero.NewOsFs()
			_, err := homedir.Dir()
			assert.NoError(t, err)
			cmd.SetOut(&buf)
			wd, _ := os.Getwd()
			fmt.Println(wd, "===========wd===========")
			afero.WriteFile(FS, "../../.metana.yml", []byte("dir: testdata\nstore: ''"), 0644)
			return RunUp(cmd, []string{}, FS, "../..")
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
