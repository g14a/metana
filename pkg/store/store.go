package store

import "github.com/go-pg/pg/v10"

type Store interface {
	Set(lastMigration LastMigration) error
	Load() (LastMigration, error)
}

type LastMigration struct {
	TimeStamp string   `json:"time_stamp" bson:"time_stamp" pg:"time_stamp"`
	tableName struct{} `pg:"migrations"`
}

type PostgresStore struct {
	db *pg.DB
}

func (p PostgresStore) Set(lastMigration LastMigration) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Close()

	_, err = tx.Model(LastMigration{}).Exec(`TRUNCATE migrations`)
	if err != nil {
		return err
	}

	_, err = tx.Model(LastMigration{}).Insert(lastMigration)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	return nil
}

func (p PostgresStore) Load() (LastMigration, error) {
	var lm LastMigration
	err := p.db.Model(LastMigration{}).Select(&lm)
	if err != nil {
		return lm, err
	}
	return lm, nil
}
