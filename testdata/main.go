// This file is auto generated. DO NOT EDIT!
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/g14a/metana/testdata/scripts"
)

func MigrateUp(upUntil string, lastRunTS int) (err error) {

	initSchemaMigration := &scripts.InitSchemaMigration{}
	initSchemaMigration.Timestamp = 1621746399
	initSchemaMigration.Filename = "1621746399-InitSchema.go"
	initSchemaMigration.MigrationName = "InitSchema"

	if lastRunTS < initSchemaMigration.Timestamp {
		fmt.Printf("  >>> Migrating up: %s\n", initSchemaMigration.Filename)
		errInitSchema := initSchemaMigration.Up()

		if errInitSchema != nil {
			fmt.Errorf("%w", errInitSchema)
		}
		fmt.Fprintf(os.Stderr, "1621746399-InitSchema.go\n")
	}

	if upUntil == "InitSchema" {
		if lastRunTS == initSchemaMigration.Timestamp {
			return
		}
		fmt.Printf("  >>> Migrated up until: %s\n", initSchemaMigration.Filename)
		return
	}

	addIndexesMigration := &scripts.AddIndexesMigration{}
	addIndexesMigration.Timestamp = 1621746406
	addIndexesMigration.Filename = "1621746406-AddIndexes.go"
	addIndexesMigration.MigrationName = "AddIndexes"

	if lastRunTS < addIndexesMigration.Timestamp {
		fmt.Printf("  >>> Migrating up: %s\n", addIndexesMigration.Filename)
		errAddIndexes := addIndexesMigration.Up()

		if errAddIndexes != nil {
			fmt.Errorf("%w", errAddIndexes)
		}
		fmt.Fprintf(os.Stderr, "1621746406-AddIndexes.go\n")
	}

	if upUntil == "AddIndexes" {
		if lastRunTS == addIndexesMigration.Timestamp {
			return
		}
		fmt.Printf("  >>> Migrated up until: %s\n", addIndexesMigration.Filename)
		return
	}

	addDataMigration := &scripts.AddDataMigration{}
	addDataMigration.Timestamp = 1621746410
	addDataMigration.Filename = "1621746410-AddData.go"
	addDataMigration.MigrationName = "AddData"

	if lastRunTS < addDataMigration.Timestamp {
		fmt.Printf("  >>> Migrating up: %s\n", addDataMigration.Filename)
		errAddData := addDataMigration.Up()

		if errAddData != nil {
			fmt.Errorf("%w", errAddData)
		}
		fmt.Fprintf(os.Stderr, "1621746410-AddData.go\n")
	}

	if upUntil == "AddData" {
		if lastRunTS == addDataMigration.Timestamp {
			return
		}
		fmt.Printf("  >>> Migrated up until: %s\n", addDataMigration.Filename)
		return
	}

	return nil

}

func MigrateDown(downUntil string, lastRunTS int) (err error) {

	addDataMigration := &scripts.AddDataMigration{}
	addDataMigration.Timestamp = 1621746410
	addDataMigration.Filename = "1621746410-AddData.go"
	addDataMigration.MigrationName = "AddData"

	if lastRunTS >= addDataMigration.Timestamp {
		fmt.Printf("  >>> Migrating down: %s\n", addDataMigration.Filename)
		errAddData := addDataMigration.Down()

		if errAddData != nil {
			fmt.Errorf("%w", errAddData)
		}
		fmt.Fprintf(os.Stderr, "1621746410-AddData.go\n")
	}

	if downUntil == "AddData" {
		if lastRunTS == addDataMigration.Timestamp {
			return
		}
		fmt.Printf("  >>> Migrated down until: %s\n", addDataMigration.Filename)
		return
	}

	addIndexesMigration := &scripts.AddIndexesMigration{}
	addIndexesMigration.Timestamp = 1621746406
	addIndexesMigration.Filename = "1621746406-AddIndexes.go"
	addIndexesMigration.MigrationName = "AddIndexes"

	if lastRunTS >= addIndexesMigration.Timestamp {
		fmt.Printf("  >>> Migrating down: %s\n", addIndexesMigration.Filename)
		errAddIndexes := addIndexesMigration.Down()

		if errAddIndexes != nil {
			fmt.Errorf("%w", errAddIndexes)
		}
		fmt.Fprintf(os.Stderr, "1621746406-AddIndexes.go\n")
	}

	if downUntil == "AddIndexes" {
		if lastRunTS == addIndexesMigration.Timestamp {
			return
		}
		fmt.Printf("  >>> Migrated down until: %s\n", addIndexesMigration.Filename)
		return
	}

	initSchemaMigration := &scripts.InitSchemaMigration{}
	initSchemaMigration.Timestamp = 1621746399
	initSchemaMigration.Filename = "1621746399-InitSchema.go"
	initSchemaMigration.MigrationName = "InitSchema"

	if lastRunTS >= initSchemaMigration.Timestamp {
		fmt.Printf("  >>> Migrating down: %s\n", initSchemaMigration.Filename)
		errInitSchema := initSchemaMigration.Down()

		if errInitSchema != nil {
			fmt.Errorf("%w", errInitSchema)
		}
		fmt.Fprintf(os.Stderr, "1621746399-InitSchema.go\n")
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
