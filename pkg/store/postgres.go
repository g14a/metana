package store

import (
	"github.com/g14a/go-migrate/pkg/types"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type PGDB struct {
	db *pg.DB
}

func (p PGDB) Set(track types.Track) error {
	err := p.CreateTable()
	if err != nil {
		return err
	}
	_, err = p.db.Model(&track).Exec(`TRUNCATE migrations`)
	if err != nil {
		return err
	}

	_, err = p.db.Model(&track).Insert(&track)
	if err != nil {
		return err
	}
	return nil
}

func (p PGDB) Load() (types.Track, error) {
	var track types.Track

	err := p.db.Model(&track).Select(&track)
	if err == pg.ErrNoRows {
		return types.Track{}, nil
	}

	return track, nil
}

func (p PGDB) CreateTable() error {
	err := p.db.Model(&types.Track{}).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})

	if err != nil {
		return err
	}

	return nil
}
