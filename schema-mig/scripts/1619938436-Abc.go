package scripts

import (
	"fmt"
)

type AbcMigration struct {
	Timestamp     int
	Filename      string
	MigrationName string
}

func (r *AbcMigration) Up() error {
	fmt.Println("Abc up")
	return nil
}

func (r *AbcMigration) Down() error {
	fmt.Println("Abc down")
	return nil
}
