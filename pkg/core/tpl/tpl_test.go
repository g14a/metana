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
	
	func ({{ .FirstChar }} *{{ .MigrationName }}Migration) Up() error {
		fmt.Println("{{ .MigrationName }} up")
		return nil
	}
	
	func ({{ .FirstChar }} *{{ .MigrationName }}Migration) Down() error {
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
	
	func ({{ .FirstChar }} *{{ .MigrationName }}Migration) Up() error {fmt.Println("template up")
return nil
	}
	
	func ({{ .FirstChar }} *{{ .MigrationName }}Migration) Down() error {fmt.Println("template down")
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
