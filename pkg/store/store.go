package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/g14a/metana/pkg/types"
	"github.com/go-pg/pg/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mconnString "go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Store interface {
	Set(track types.Track) error
	Load() (types.Track, error)
	Wipe() error
}

func GetStoreViaConn(connString string, dir string) (Store, error) {

	if strings.HasPrefix(connString, "@") {
		connString = strings.TrimPrefix(connString, "@")
		connString = os.Getenv(connString)
	}

	switch {
	case strings.Contains(connString, "postgres://"):
		pgOptions, err := pg.ParseURL(connString)
		if err != nil {
			log.Println("Couldn't parse postgres connection string")
		}

		db := pg.Connect(pgOptions)

		_, err = db.Exec("SELECT 1")
		if err != nil {
			return nil, fmt.Errorf("Could not connect to your PostgreSQL DB, ERROR: %w", err)
		}

		p := PGDB{db: db}
		err = p.CreateTable()
		if err != nil {
			return nil, fmt.Errorf("Could not create migrations table in postgres")
		}

		return p, nil
	case strings.Contains(connString, "mongodb"):
		ctx := context.TODO()
		cs, err := mconnString.ParseAndValidate(connString)
		clientOptions := options.Client().ApplyURI(connString)
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			return nil, fmt.Errorf("Could not connec to MongoDB, ERROR: %w", err)
		}
		err = client.Ping(ctx, nil)
		if err != nil {
			return nil, err
		}
		return MongoDb{coll: *client.Database(cs.Database).Collection("migrations")}, nil
	}

	jsonFile, err := os.OpenFile(dir+"/migrate.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return File{file: *jsonFile}, nil
}

func TrackToSetDown(track types.Track, num int) types.Track {
	if len(track.Migrations) == num {
		track.LastRun = track.Migrations[len(track.Migrations)-num].Title
		track.LastRunTS = track.Migrations[len(track.Migrations)-num].Timestamp
	} else {
		track.LastRun = track.Migrations[len(track.Migrations)-num-1].Title
		track.LastRunTS = track.Migrations[len(track.Migrations)-num-1].Timestamp
	}
	track.Migrations = track.Migrations[:len(track.Migrations)-num]
	if len(track.Migrations) == 0 {
		return types.Track{}
	}

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
