// This file is auto generated. DO NOT EDIT!
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/g14a/metana/testdata/environments/dev/scripts"
)

func MigrateUp(upUntil string, lastRunTS int) (err error) {

	initSchemaMigration := &scripts.InitSchemaMigration{}
	initSchemaMigration.Timestamp = 1623502023
	initSchemaMigration.Filename = "1623502023-InitSchema.go"
	initSchemaMigration.MigrationName = "InitSchema"

	if lastRunTS < initSchemaMigration.Timestamp {
		fmt.Printf("  >>> Migrating up: %s\n", initSchemaMigration.Filename)
		errInitSchema := initSchemaMigration.Up()

		if errInitSchema != nil {
			fmt.Errorf("%w", errInitSchema)
		}
		fmt.Fprintf(os.Stderr, "1623502023-InitSchema.go\n")
	}

	if upUntil == "InitSchema" {
		if lastRunTS == initSchemaMigration.Timestamp {
			return
		}
		fmt.Printf("  >>> Migrated up until: %s\n", initSchemaMigration.Filename)
		return
	}

	return nil

}

func MigrateDown(downUntil string, lastRunTS int) (err error) {

	initSchemaMigration := &scripts.InitSchemaMigration{}
	initSchemaMigration.Timestamp = 1623502023
	initSchemaMigration.Filename = "1623502023-InitSchema.go"
	initSchemaMigration.MigrationName = "InitSchema"

	if lastRunTS >= initSchemaMigration.Timestamp {
		fmt.Printf("  >>> Migrating down: %s\n", initSchemaMigration.Filename)
		errInitSchema := initSchemaMigration.Down()

		if errInitSchema != nil {
			fmt.Errorf("%w", errInitSchema)
		}
		fmt.Fprintf(os.Stderr, "1623502023-InitSchema.go\n")
	}

	if downUntil == "InitSchema" {
		if lastRunTS == initSchemaMigration.Timestamp {
			return
		}
		fmt.Printf("  >>> Migrated down until: %s\n", initSchemaMigration.Filename)
		return
	}

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
