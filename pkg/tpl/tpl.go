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
	"os"
	"{{ .pwd }}/migrations/scripts"
	"fmt"
)

func MigrateUp() error {
	track, _ := Load()
	
	return nil
}

func MigrateDown() error {
	
	return nil
}

func main() {
	if os.Args[1] == "up" {
		err := MigrateUp()
		if err != nil {
			fmt.Println(err)
		}
	}

	if os.Args[1] == "down" {
		err := MigrateDown()
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
		fmt.Printf("  >>> Running up: %s\n\n", {{ .Lower }}Migration.Filename)
		err{{ .MigrationName }} := {{ .Lower }}Migration.Up()

		if err{{ .MigrationName }} != nil {
			return fmt.Errorf("{{ .Filename }}, %w", err{{ .MigrationName }})
		}
	
		err{{ .MigrationName }} = Set({{ .Lower }}Migration.Timestamp, "{{ .Filename }}")
		if err{{ .MigrationName }} != nil {
			return fmt.Errorf("{{ .Filename }}, %w", err{{ .MigrationName }})
		}
	}

	return nil
`)
	} else {
		return []byte(`
	{{ .Lower }}Migration := &scripts.{{ .MigrationName }}Migration{}
	{{ .Lower }}Migration.Timestamp = {{ .Timestamp }}
	{{ .Lower }}Migration.Filename = "{{ .Filename }}"
	{{ .Lower }}Migration.MigrationName = "{{ .MigrationName }}"

	err{{ .MigrationName }} := {{ .Lower }}Migration.Down()

	if err{{ .MigrationName }} != nil {
		return fmt.Errorf("{{ .Filename }}, %w", err{{ .MigrationName }})
	}

	err{{ .MigrationName }} = Set({{ .Lower }}Migration.Timestamp, "{{ .Filename }}")
	if err{{ .MigrationName }} != nil {
		return fmt.Errorf("{{ .Filename }}, %w", err{{ .MigrationName }})
	}	

	return nil
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

func Set(timestamp int, fileName string) error {
	track, err := Load()
	if err != nil {
		return err
	}

	track.LastRun = fileName
	track.LastRunTS  = timestamp
	track.Migrations = append(track.Migrations, types.Migration{
		Title:     fileName,
		Timestamp: timestamp,
	})

	bytes, err := json.MarshalIndent(track,"", "	")
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
}
