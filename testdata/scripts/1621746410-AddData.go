package scripts

import (
	"fmt"
)

type AddDataMigration struct {
	Timestamp     int
	Filename      string
	MigrationName string
}

func (r *AddDataMigration) Up() error {
	fmt.Println("AddData up")
	return nil
}

func (r *AddDataMigration) Down() error {
	fmt.Println("AddData down")
	return nil
}
