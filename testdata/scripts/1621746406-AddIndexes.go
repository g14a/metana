package scripts

import (
	"fmt"
)

type AddIndexesMigration struct {
	Timestamp     int
	Filename      string
	MigrationName string
}

func (r *AddIndexesMigration) Up() error {
	fmt.Println("AddIndexes up")
	return nil
}

func (r *AddIndexesMigration) Down() error {
	fmt.Println("AddIndexes down")
	return nil
}
