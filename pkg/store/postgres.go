package store

import (
	"github.com/g14a/go-migrate/pkg/types"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type PGDB struct {
	db *pg.DB
}

func (p PGDB) Set(timestamp int, filename string, up bool) error {
	return nil
}

func (p PGDB) Load() (types.Track, error) {
	return nil, nil
}

func (p PGDB) CreateTable() error {
	err := p.db.Model(types.Track{}).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})

	if err != nil {
		return err
	}

	return nil
}
