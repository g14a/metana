package tpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrationTemplate(t *testing.T) {
	resultTpl := MigrationTemplate("", "")

	assert.Equal(t, []byte(`package scripts

	import (
		"fmt"
	)
	
	type {{ .MigrationName }}Migration struct {
		Timestamp int
		Filename  string
		MigrationName string
	}
	
	func (r *{{ .MigrationName }}Migration) Up() error {
		fmt.Println("{{ .MigrationName }} up")
		return nil
	}
	
	func (r *{{ .MigrationName }}Migration) Down() error {
		fmt.Println("{{ .MigrationName }} down")
		return nil
	}
`), resultTpl)
}

func TestMigrationTemplate_Custom(t *testing.T) {
	resultTpl := MigrationTemplate(`fmt.Println("template up")`+"\nreturn nil", `fmt.Println("template down")`+"\nreturn nil")

	assert.Equal(t, `package scripts

	import (
		"fmt"
	)
	
	type {{ .MigrationName }}Migration struct {
		Timestamp int
		Filename  string
		MigrationName string
	}
	
	func (r *{{ .MigrationName }}Migration) Up() error {fmt.Println("template up")
return nil
	}
	
	func (r *{{ .MigrationName }}Migration) Down() error {fmt.Println("template down")
return nil
	}
`, string(resultTpl))
}

func TestInitMigrationRunTemplate(t *testing.T) {
	resultTpl := InitMigrationRunTemplate()

	assert.Equal(t, []byte(`// This file is auto generated. DO NOT EDIT!
	package main
	
	import (
		"flag"
		"fmt"
		"os"

		"{{ .pwd }}/{{ .dir }}/scripts"
	)
	
	func MigrateUp(upUntil string, lastRunTS int) (err error) {
		return nil
	}
	
	func MigrateDown(downUntil string, lastRunTS int) (err error) {
		return nil
	}
	
	func main() {
		upCmd := flag.NewFlagSet("up", flag.ExitOnError)
		downCmd := flag.NewFlagSet("down", flag.ExitOnError)
	
		var upUntil, downUntil string
		var lastRunTS int
		upCmd.StringVar(&upUntil, "until", "", "")
		upCmd.IntVar(&lastRunTS, "last-run-ts", 0, "")
		downCmd.StringVar(&downUntil, "until", "", "")
		downCmd.IntVar(&lastRunTS, "last-run-ts", 0, "")
	
		switch os.Args[1] {
		case "up":
			err := upCmd.Parse(os.Args[2:])
			if err != nil {
				return
			}
		case "down":
			err := downCmd.Parse(os.Args[2:])
			if err != nil {
				return
			}
		}
	
		if upCmd.Parsed() {
			err := MigrateUp(upUntil, lastRunTS)
			if err != nil {
				fmt.Fprintf(os.Stdout, err.Error())
			}
		}
	
		if downCmd.Parsed() {
			err := MigrateDown(downUntil, lastRunTS)
			if err != nil {
				fmt.Fprintf(os.Stdout, err.Error())
			}
		}
	}
`), resultTpl)
}

func TestAddMigrationTemplate(t *testing.T) {
	tests := []struct {
		up        bool
		wantedTpl string
	}{
		{
			up: true,
			wantedTpl: `
	{{ .Lower }}Migration := &scripts.{{ .MigrationName }}Migration{}
	{{ .Lower }}Migration.Timestamp = {{ .Timestamp }}
	{{ .Lower }}Migration.Filename = "{{ .Filename }}"
	{{ .Lower }}Migration.MigrationName = "{{ .MigrationName }}"

	if lastRunTS < {{ .Lower }}Migration.Timestamp {
		fmt.Printf("  >>> Migrating up: %s\n", {{ .Lower }}Migration.Filename)
		err{{ .MigrationName }} := {{ .Lower }}Migration.Up()

		if err{{ .MigrationName }} != nil {
			fmt.Errorf("%w", err{{ .MigrationName }})
		}
		fmt.Fprintf(os.Stderr, "{{ .Filename }}\n")
	}

	if upUntil == "{{ .MigrationName }}" {
		if lastRunTS == {{ .Lower }}Migration.Timestamp {
			return
		}
		fmt.Printf("  >>> Migrated up until: %s\n", {{ .Lower }}Migration.Filename)
		return
	}

	return nil
`,
		},
		{
			up: false,
			wantedTpl: `
	{{ .Lower }}Migration := &scripts.{{ .MigrationName }}Migration{}
	{{ .Lower }}Migration.Timestamp = {{ .Timestamp }}
	{{ .Lower }}Migration.Filename = "{{ .Filename }}"
	{{ .Lower }}Migration.MigrationName = "{{ .MigrationName }}"

	if lastRunTS >= {{ .Lower }}Migration.Timestamp {
		fmt.Printf("  >>> Migrating down: %s\n", {{ .Lower }}Migration.Filename)
		err{{ .MigrationName }} := {{ .Lower }}Migration.Down()

		if err{{ .MigrationName }} != nil {
			fmt.Errorf("%w", err{{ .MigrationName }})
		}
		fmt.Fprintf(os.Stderr, "{{ .Filename }}\n")
	}

	if downUntil == "{{ .MigrationName }}" {
		if lastRunTS == {{ .Lower }}Migration.Timestamp {
			return
		}
		fmt.Printf("  >>> Migrated down until: %s\n", {{ .Lower }}Migration.Filename)
		return
	}
`,
		},
	}

	for _, tt := range tests {
		resultTpl := AddMigrationTemplate(tt.up)
		assert.Equal(t, tt.wantedTpl, string(resultTpl))
	}
}
