package tpl

func MigrationTemplate() []byte {
	return []byte(`package scripts

import (
	"fmt"
)

type {{ .MigrationName }}Migration struct {
	Timestamp int
	Filename  string
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
	{{ .Lower }}Migration.Filename = {{ .Filename }}
	{{ .Lower }}Migration.MigrationName = {{ .MigrationName }}

	err{{ .MigrationName }} := {{ .Lower }}Migration.Up()

	if err{{ .MigrationName }} != nil {
		return fmt.Errorf("{{ .Filename }}, %w", err{{ .MigrationName }})
	}

	err{{ .MigrationName }} = Set({{ .Lower }}Migration.Timestamp, "{{ .Filename }}")
	if err{{ .MigrationName }} != nil {
		return fmt.Errorf("{{ .Filename }}, %w", err{{ .MigrationName }})
	}

	return nil
`)
	} else {
		return []byte(`
	{{ .Lower }}Migration := &scripts.{{ .MigrationName }}Migration{}
	{{ .Lower }}Migration.Timestamp = {{ .Timestamp }}
	{{ .Lower }}Migration.Filename = {{ .Filename }}

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
)

func Set(timestamp int, fileName string) error {
	track, err := Load()
	if err != nil {
		return err
	}

	track.LastRun = fileName
	track.Migrations = append(track.Migrations, Migration{
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

func Load() (Track, error) {
	track, err := os.ReadFile("migrate.json")
	if err != nil {
		return Track{}, err
	}

	t := Track{}

	if len(track) > 0 {
		err = json.Unmarshal(track, &t)
		if err != nil {
			return Track{}, err
		}
	}

	return t, nil
}

type Track struct {
	LastRun    string      
	Migrations []Migration 
}

type Migration struct {
	Title     string 
	Timestamp int   
}`)
}

type NewMigration struct {
	Lower         string
	MigrationName string
	Timestamp     string
	Filename      string
}
