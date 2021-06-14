package cmd

import (
	"bytes"
	"testing"
	"time"

	"github.com/g14a/metana/pkg"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_List(t *testing.T) {
	var buf bytes.Buffer
	metanaCmd := NewMetanaCommand()

	listCmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			FS := afero.NewMemMapFs()
			cmd.SetOut(&buf)
			FS.MkdirAll("/Users/g14a/metana/migrations/environments/dev/scripts", 0755)
			afero.WriteFile(FS, "/Users/g14a/metana/migrations/environments/dev/scripts/1621081055-InitSchema.go", []byte("{}"), 0644)
			afero.WriteFile(FS, "/Users/g14a/metana/migrations/environments/dev/scripts/1621084125-AddIndexes.go", []byte("{}"), 0644)
			afero.WriteFile(FS, "/Users/g14a/metana/migrations/environments/dev/scripts/1621084135-AddFKeys.go", []byte("{}"), 0644)
			return RunList(cmd, "/Users/g14a/metana", FS)
		},
	}
	listCmd.Flags().StringP("dir", "d", "", "Specify migrations dir")
	listCmd.Flags().StringP("env", "e", "", "List migrations in an environment")

	metanaCmd.AddCommand(listCmd)
	_, err := pkg.ExecuteCommand(metanaCmd, "list", "--env=dev")
	assert.NoError(t, err)

	var data [][]string

	data = append(data, []string{"1621081055-InitSchema.go", time.Now().Format("02-01-2006 15:04")})
	data = append(data, []string{"1621084125-AddIndexes.go", time.Now().Format("02-01-2006 15:04")})
	data = append(data, []string{"1621084135-AddFKeys.go", time.Now().Format("02-01-2006 15:04")})

	var outBuf bytes.Buffer
	table := tablewriter.NewWriter(&outBuf)
	table.SetHeader([]string{"Migration", "Last Modified"})

	for _, row := range data {
		table.Append(row)
	}
	table.Render()

	assert.Equal(t, len(outBuf.String()), len(buf.String()))
}
