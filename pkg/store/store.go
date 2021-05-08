package store

import (
	"github.com/g14a/go-migrate/pkg/types"
	"github.com/go-pg/pg/v10"
	"log"
	"strconv"
	"strings"
)

type Store interface {
	Set(track types.Track) error
	Load() (types.Track, error)
}

func GetStoreViaConn(connString string) Store {
	//switch {
	//case strings.Contains(connString, "postgres://"):
	//	options, err := pg.ParseURL(connString)
	//	if err != nil {
	//		log.Println("Couldn't parse postgres connection string")
	//	}
	//
	//	db := pg.Connect(options)
	//
	//	_, err = db.Exec("SELECT 1")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	p := PGDB{db: db}
	//	err = p.CreateTable()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	return p
	//}

	options, err := pg.ParseURL(connString)
	if err != nil {
		log.Println("Couldn't parse postgres connection string")
	}

	db := pg.Connect(options)

	_, err = db.Exec("SELECT 1")
	if err != nil {
		log.Fatal(err)
	}

	p := PGDB{db: db}
	err = p.CreateTable()
	if err != nil {
		log.Fatal(err)
	}
	return p
}

func TrackToSet(track types.Track, timestamp int, filename string, up bool) (types.Track, error) {
	if up {
		track.LastRun = filename
		track.LastRunTS = timestamp
		track.Migrations = append(track.Migrations, types.Migration{
			Title:     filename,
			Timestamp: timestamp,
		})
	} else {
		track.Migrations = track.Migrations[:len(track.Migrations)-1]
		if len(track.Migrations) == 0 {
			return types.Track{}, nil
		}
		track.LastRun = track.Migrations[len(track.Migrations)-1].Title
		track.LastRunTS = track.Migrations[len(track.Migrations)-1].Timestamp
	}

	return track, nil
}

func ProcessLogs(logs string) types.Track {
	track := types.Track{}
	lines := strings.Split(logs,"\n")
	for _, line := range lines {
		if len(line) > 0 {
			migration := types.Migration{}
			track.LastRun = line
			migration.Title = line
			line = strings.TrimSuffix(line, ".go")
			migArr := strings.Split(line, "-")
			timestamp, err := strconv.Atoi(migArr[0])
			if err != nil {
				log.Fatal(err)
			}
			track.LastRunTS = timestamp
			migration.Timestamp = timestamp
			track.Migrations = append(track.Migrations, migration)
		}
	}

	return track
}