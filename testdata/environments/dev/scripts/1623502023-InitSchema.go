package scripts

import (
	"fmt"
)

type InitSchemaMigration struct {
	Timestamp     int
	Filename      string
	MigrationName string
}

func (i *InitSchemaMigration) Up() error {
	fmt.Println("InitSchema up")
	return nil
}

func (i *InitSchemaMigration) Down() error {
	fmt.Println("InitSchema down")
	return nil
}
