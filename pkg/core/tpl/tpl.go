package tpl

func MigrationTemplate(customUp, customDown string) []byte {
	finalUpComponent := `
		fmt.Println("{{ .MigrationName }} up")
		return nil`

	if customUp != "" {
		finalUpComponent = customUp
	}

	finalDownComponent := `
		fmt.Println("{{ .MigrationName }} down")
		return nil`

	if customDown != "" {
		finalDownComponent = customDown
	}

	return []byte(`package scripts

	import (
		"fmt"
	)
	
	type {{ .MigrationName }}Migration struct {
		Timestamp int
		Filename  string
		MigrationName string
	}
	
	func ({{ .FirstChar }} *{{ .MigrationName }}Migration) Up() error {` +
		finalUpComponent + `
	}
	
	func ({{ .FirstChar }} *{{ .MigrationName }}Migration) Down() error {` +
		finalDownComponent + `
	}
`)
}

func InitMigrationRunTemplate() []byte {
	return []byte(`// This file is auto generated. DO NOT EDIT!
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
`)
}

func NewAddMigrationTemplate() []byte {
	return []byte(`// This file is auto generated. DO NOT EDIT!
package main

import (
	"flag"
	"fmt"
	"os"

	"{{ .Pwd }}/{{ .Dir }}/scripts"
)

func MigrateUp(upUntil string, lastRunTS int) (err error) {
        {{ range $m := .Create }}
        {{ $m.Lower }}Migration := &scripts.{{ $m.MigrationName }}Migration{}
        	{{ $m.Lower }}Migration.Timestamp = {{ $m.Timestamp }}
        	{{ $m.Lower }}Migration.Filename = "{{ $m.Filename }}"
        	{{ $m.Lower }}Migration.MigrationName = "{{ $m.MigrationName }}"

        	if lastRunTS < {{ $m.Lower }}Migration.Timestamp {
        		fmt.Printf("  >>> Migrating up: %s\n", {{ $m.Lower }}Migration.Filename)
        		err{{ $m.MigrationName }} := {{ $m.Lower }}Migration.Up()

        		if err{{ $m.MigrationName }} != nil {
        			fmt.Errorf("%w", err{{ $m.MigrationName }})
        		}
        		fmt.Fprintf(os.Stderr, "{{ $m.Filename }}\n")
        	}

        	if upUntil == "{{ $m.MigrationName }}" {
        		if lastRunTS == {{ $m.Lower }}Migration.Timestamp {
        			return
        		}
        		fmt.Printf("  >>> Migrated up until: %s\n", {{ $m.Lower }}Migration.Filename)
        		return
        	}
		{{ end }}
		return nil
}

func MigrateDown(downUntil string, lastRunTS int) (err error) {
		{{ range $m := .Create }}
        {{ $m.Lower }}Migration := &scripts.{{ $m.MigrationName }}Migration{}
		{{ $m.Lower }}Migration.Timestamp = {{ $m.Timestamp }}
		{{ $m.Lower }}Migration.Filename = "{{ $m.Filename }}"
		{{ $m.Lower }}Migration.MigrationName = "{{ $m.MigrationName }}"

		if lastRunTS >= {{ $m.Lower }}Migration.Timestamp {
			fmt.Printf("  >>> Migrating down: %s\n", {{ $m.Lower }}Migration.Filename)
			err{{ $m.MigrationName }} := {{ $m.Lower }}Migration.Down()

			if err{{ $m.MigrationName }} != nil {
				fmt.Errorf("%w", err{{ $m.MigrationName }})
			}
			fmt.Fprintf(os.Stderr, "{{ $m.Filename }}\n")
		}

		if downUntil == "{{ $m.MigrationName }}" {
			if lastRunTS == {{ $m.Lower }}Migration.Timestamp {
				return
			}
			fmt.Printf("  >>> Migrated down until: %s\n", {{ $m.Lower }}Migration.Filename)
			return
		}
		{{ end }}
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
}`)
}

type NewMigration struct {
	Lower         string
	MigrationName string
	Timestamp     string
	Filename      string
	MigrationsDir string
	FirstChar     string
}
