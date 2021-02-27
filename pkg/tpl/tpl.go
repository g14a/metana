package tpl

func MigrationTemplate() []byte {
	return []byte(`package main

import (
	"fmt"
)

type {{ .Name }}Migration struct {

}

func (r *{{ .Name }}Migration) Up()  {
	fmt.Println("migration up")
}

func (r *{{ .Name }}Migration) Down()  {
	fmt.Println("migration down")
}
`)
}

func InitMigrationRunTemplate() []byte {
	return []byte(`package main

import (
	"os"
)

func MigrateUp() {
	
}

func MigrateDown() {

}

func main() {
	if os.Args[0] == "up" {
		MigrateUp()
	}

	if os.Args[0] == "down" {
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
