package store

import (
	"github.com/g14a/go-migrate/pkg/types"
	"github.com/go-pg/pg/v10"
	"log"
	"os"
	"strconv"
	"strings"
)

type Store interface {
	Set(track types.Track) error
	Load() (types.Track, error)
}

func GetStoreViaConn(connString string, dir string) Store {

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

		p := PGDB{db: db}
		err = p.CreateTable()
		if err != nil {
			log.Fatal(err)
		}

		return p
	}

	jsonFile, err := os.OpenFile(dir+"/migrate.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return File{file: *jsonFile}
}

func TrackToSetDown(track types.Track, num int) types.Track {

	track.LastRun = track.Migrations[len(track.Migrations)-num-1].Title
	track.LastRunTS = track.Migrations[len(track.Migrations)-num-1].Timestamp
	track.Migrations = track.Migrations[:len(track.Migrations)-num]

	return track
}

func ProcessLogs(logs string) (types.Track, int) {
	track := types.Track{}
	lines := strings.Split(logs, "\n")
	num := 0
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
			num++
		}
	}

	return track, num
}
