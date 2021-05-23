package scripts

import (
	"fmt"
)

type InitSchemaMigration struct {
	Timestamp     int
	Filename      string
	MigrationName string
}

func (r *InitSchemaMigration) Up() error {
	fmt.Println("InitSchema up")
	return nil
}

func (r *InitSchemaMigration) Down() error {
	fmt.Println("InitSchema down")
	return nil
}
