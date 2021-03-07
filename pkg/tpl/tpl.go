package tpl

func MigrationTemplate() []byte {
	return []byte(`package main

import (
	"fmt"
)

type {{ .Name }}Migration struct {
}

func (r *{{ .Name }}Migration) Up() error {
	fmt.Println("{{ .Name }} up")
	return nil
}

func (r *{{ .Name }}Migration) Down() error {
	fmt.Println("{{ .Name }} down")
	return nil
}
`)
}

func InitMigrationRunTemplate() []byte {
	return []byte(`// This file is auto generated. DO NOT EDIT!
package main

import (
	"os"
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

type NewMigration struct {
	Name      string
	Timestamp string
}
