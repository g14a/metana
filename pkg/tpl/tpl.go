package tpl

func MigrationTemplate() []byte {
	return []byte(`package main

import (
	"fmt"
)

type {{ .Name }}Migration struct {

}

func (r *{{ .Name }}Migration) Up()  {
	fmt.Println("{{ .Name }} up")
}

func (r *{{ .Name }}Migration) Down()  {
	fmt.Println("{{ .Name }} down")
}
`)
}

func InitMigrationRunTemplate() []byte {
	return []byte(`// This file is auto generated. DO NOT EDIT!
package main

import (
	_interface "go-migrate/migrations/interfaces"
	"os"
)

func MigrateUp() {
	
}

func MigrateDown() {

}

func main() {
	if os.Args[1] == "up" {
		MigrateUp()
	}

	if os.Args[1] == "down" {
		MigrateDown()
	}
}
`)
}

func InitMigrationTemplate() []byte {
	return []byte(`package _interface

type Migration interface {
	Up()
	Down()
}
`)
}

type NewMigration struct {
	Name      string
	Timestamp string
}
