package tpl

func MigrationTemplate() []byte {
	return []byte(`package scripts

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
`)
}

func InitMigrationRunTemplate() []byte {
	return []byte(`// This file is auto generated. DO NOT EDIT!
package main

import (
	"flag"
	"os"
	"{{ .pwd }}/{{ .dir }}/scripts"
	"fmt"
)

func MigrateUp(upUntil string) (err error) {
	track, _ := Load()
	
	return nil
}

func MigrateDown(downUntil string) (err error) {
	track, _ := Load()

	return nil
}

func main() {
	upCmd := flag.NewFlagSet("up", flag.ExitOnError)
	downCmd := flag.NewFlagSet("down", flag.ExitOnError)

	var upUntil, downUntil string
	upCmd.StringVar(&upUntil, "until", "", "")
	downCmd.StringVar(&downUntil, "until", "", "")

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
		err := MigrateUp(upUntil)
		if err != nil {
			fmt.Println(err)
		}
	}

	if downCmd.Parsed() {
		err := MigrateDown(downUntil)
		if err != nil {
			fmt.Println(err)
		}
	}
}
`)
}

func AddMigrationTemplate(up bool) []byte {
	if up {
		return []byte(`
	{{ .Lower }}Migration := &scripts.{{ .MigrationName }}Migration{}
	{{ .Lower }}Migration.Timestamp = {{ .Timestamp }}
	{{ .Lower }}Migration.Filename = "{{ .Filename }}"
	{{ .Lower }}Migration.MigrationName = "{{ .MigrationName }}"

	if track.LastRunTS < {{ .Lower }}Migration.Timestamp {
		fmt.Printf("\n  >>> Migrating up: %s\n", {{ .Lower }}Migration.Filename)
		err{{ .MigrationName }} := {{ .Lower }}Migration.Up()

		if err{{ .MigrationName }} != nil {
			return fmt.Errorf("{{ .Filename }}, %w", err{{ .MigrationName }})
		}
	
		err{{ .MigrationName }} = Set({{ .Lower }}Migration.Timestamp, {{ .Lower }}Migration.Filename, true)
		if err{{ .MigrationName }} != nil {
			return fmt.Errorf("{{ .Filename }}, %w", err{{ .MigrationName }})
		}
	}

	if upUntil == "{{ .MigrationName }}" {
		if track.LastRunTS == {{ .Lower }}Migration.Timestamp {
			return
		}
		fmt.Printf("\n  >>> Migrated up until: %s\n", {{ .Lower }}Migration.Filename)
		return
	}

	return nil
`)
	} else {
		return []byte(`
	{{ .Lower }}Migration := &scripts.{{ .MigrationName }}Migration{}
	{{ .Lower }}Migration.Timestamp = {{ .Timestamp }}
	{{ .Lower }}Migration.Filename = "{{ .Filename }}"
	{{ .Lower }}Migration.MigrationName = "{{ .MigrationName }}"

	if track.LastRunTS >= {{ .Lower }}Migration.Timestamp {
		fmt.Printf("\n  >>> Migrating down: %s\n", {{ .Lower }}Migration.Filename)
		err{{ .MigrationName }} := {{ .Lower }}Migration.Down()

		if err{{ .MigrationName }} != nil {
			return fmt.Errorf("{{ .Filename }}, %w", err{{ .MigrationName }})
		}

		err{{ .MigrationName }} = Set({{ .Lower }}Migration.Timestamp, {{ .Lower }}Migration.Filename, false)
		if err{{ .MigrationName }} != nil {
			return fmt.Errorf("{{ .Filename }}, %w", err{{ .MigrationName }})
		}	
	}

	if downUntil == "{{ .MigrationName }}" {
		if track.LastRunTS == {{ .Lower }}Migration.Timestamp {
			return
		}
		fmt.Printf("  >>> Migrated down until: %s\n", {{ .Lower }}Migration.Filename)
		return
	}
`)
	}
}

func StoreTemplate() []byte {
	return []byte(`package main

import (
	"encoding/json"
	"os"

	"github.com/g14a/go-migrate/pkg/types"
)

func Set(timestamp int, fileName string, up bool) error {
	track, err := Load()
	if err != nil {
		return err
	}

	if up {
		track.LastRun = fileName
		track.LastRunTS = timestamp
		track.Migrations = append(track.Migrations, types.Migration{
			Title:     fileName,
			Timestamp: timestamp,
		})
	} else {
		if len(track.Migrations) == 0 {
			err = os.WriteFile("migrate.json", nil, 0644)
			if err != nil {
				return err
			}
			return nil
		}
		track.LastRun = fileName
		track.LastRunTS = timestamp
		track.Migrations = track.Migrations[:len(track.Migrations)-1]
	}

	bytes, err := json.MarshalIndent(track, "", "	")
	if err != nil {
		return err
	}

	err = os.WriteFile("migrate.json", bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func Load() (types.Track, error) {
	track, err := os.ReadFile("migrate.json")
	if err != nil {
		return types.Track{}, err
	}

	t := types.Track{}

	if len(track) > 0 {
		err = json.Unmarshal(track, &t)
		if err != nil {
			return types.Track{}, err
		}
	}

	return t, nil
}`)
}

type NewMigration struct {
	Lower         string
	MigrationName string
	Timestamp     string
	Filename      string
	MigrationsDir string
}
