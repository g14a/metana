package store

import (
	"fmt"
	"github.com/g14a/go-migrate/pkg/types"
	"github.com/go-pg/pg/v10"
	"log"
	"strings"
)

type Store interface {
	Set(timestamp int, fileName string, up bool) error
	Load() (types.Track, error)
}

func GetStoreViaConn(connString string) Store {
	switch {
	case strings.Contains(connString, "postgres://"):
		options, err := pg.ParseURL(connString)
		if err != nil {
			log.Println("Couldn't parse postgres connection string")
		}

		db := pg.Connect(options)

		_, err = db.Exec("SELECT 1")
		if err != nil {
			log.Fatal(err)
		}

		return PGDB{db: db}
	}

	fmt.Println("could not identify store")
	return Store
}
