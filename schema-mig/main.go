// This file is auto generated. DO NOT EDIT!
package main

import (
	"flag"
	"fmt"
	"github.com/g14a/go-migrate/schema-mig/scripts"
	"os"
)

func MigrateUp(upUntil string) (err error) {
	track, _ := Load()

	abcMigration := &scripts.AbcMigration{}
	abcMigration.Timestamp = 1619938436
	abcMigration.Filename = "1619938436-Abc.go"
	abcMigration.MigrationName = "Abc"

	if track.LastRunTS < abcMigration.Timestamp {
		fmt.Printf("\n  >>> Migrating up: %s\n", abcMigration.Filename)
		errAbc := abcMigration.Up()

		if errAbc != nil {
			return fmt.Errorf("1619938436-Abc.go, %w", errAbc)
		}

		errAbc = Set(abcMigration.Timestamp, abcMigration.Filename, true)
		if errAbc != nil {
			return fmt.Errorf("1619938436-Abc.go, %w", errAbc)
		}
	}

	if upUntil == "Abc" {
		if track.LastRunTS == abcMigration.Timestamp {
			return
		}
		fmt.Printf("\n  >>> Migrated up until: %s\n", abcMigration.Filename)
		return
	}

	return nil

}

func MigrateDown(downUntil string) (err error) {
	track, _ := Load()

	abcMigration := &scripts.AbcMigration{}
	abcMigration.Timestamp = 1619938436
	abcMigration.Filename = "1619938436-Abc.go"
	abcMigration.MigrationName = "Abc"

	if track.LastRunTS >= abcMigration.Timestamp {
		fmt.Printf("\n  >>> Migrating down: %s\n", abcMigration.Filename)
		errAbc := abcMigration.Down()

		if errAbc != nil {
			return fmt.Errorf("1619938436-Abc.go, %w", errAbc)
		}

		errAbc = Set(abcMigration.Timestamp, abcMigration.Filename, false)
		if errAbc != nil {
			return fmt.Errorf("1619938436-Abc.go, %w", errAbc)
		}
	}

	if downUntil == "Abc" {
		if track.LastRunTS == abcMigration.Timestamp {
			return
		}
		fmt.Printf("  >>> Migrated down until: %s\n", abcMigration.Filename)
		return
	}

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
