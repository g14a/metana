package store

import (
	"github.com/g14a/metana/pkg/types"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/spf13/afero"
)

type PGDB struct {
	db *pg.DB
}

func (p PGDB) Set(track types.Track, FS afero.Fs) error {
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

func (p PGDB) Load(FS afero.Fs) (types.Track, error) {
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

func (p PGDB) Wipe(FS afero.Fs) error {
	_, err := p.db.Model(&types.Track{}).Exec(`TRUNCATE migrations`)
	if err != nil {
		return err
	}

	return nil
}
