package wipe

import (
	s "github.com/g14a/metana/pkg/store"
	"os"
)

func Wipe(migrationsDir string, storeConn string) error {
	store, err := s.GetStoreViaConn(storeConn, migrationsDir)
	if err != nil {
		return err
	}

	track, err := store.Load()
	if err != nil {
		return err
	}


	for _, m := range track.Migrations {
		err := os.Remove(migrationsDir + "/scripts/" + m.Title)
		if err != nil {
			return err
		}
	}

	return nil
}
