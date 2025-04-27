package tpl

// StandaloneMigrationTemplate returns the Go source code for a standalone migration file.
// This file is compiled and executed during migrations using `go run`.
// The output must include "__COMPLETE__: <filename>" on success for Metana to track it.
func StandaloneMigrationTemplate() []byte {
	upBody := `fmt.Println("{{ .MigrationName }} up")`
	downBody := `fmt.Println("{{ .MigrationName }} down")`

	return []byte(`//go:build ignore
// +build ignore

// ⚠️ AUTO-GENERATED FILE. DO NOT IMPORT THIS INTO YOUR MAIN APPLICATION.
// ⚠️ ONLY modify the 'up()' and 'down()' functions to write your migration logic.

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
)

// ✅ UP migration logic.
// If this returns an error, Metana will auto-trigger the DOWN rollback for this migration.
func up() error {
	` + upBody + `
	return nil
}

// ✅ DOWN (rollback) logic.
// This will be triggered if 'up()' fails during execution.
func down() error {
	` + downBody + `
	return nil
}

// 🚫 DO NOT MODIFY.
// Handles flag parsing, error propagation, and execution tracking.
func main() {
	mode := flag.String("mode", "up", "migration mode: up or down")
	flag.Parse()

	var err error
	switch *mode {
	case "up":
		err = up()
	case "down":
		err = down()
	default:
		fmt.Fprintln(os.Stderr, "invalid mode: must be 'up' or 'down'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		debug.PrintStack()
		os.Exit(1)
	}

	fmt.Printf("__COMPLETE__[%s]: %s\n", *mode, "{{ .Filename }}")
}
`)
}
