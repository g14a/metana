package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/afero"

	"github.com/g14a/metana/pkg/types"
	"github.com/go-pg/pg/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mconnString "go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Store interface {
	Set(track types.Track, FS afero.Fs) error
	Load(FS afero.Fs) (types.Track, error)
}

func GetStoreViaConn(connString string, dir string, FS afero.Fs, wd string) (Store, error) {

	if strings.HasPrefix(connString, "@") {
		connString = strings.TrimPrefix(connString, "@")
		connString = os.Getenv(connString)
	}

	switch {
	case strings.Contains(connString, "postgres://"):
		pgOptions, err := pg.ParseURL(connString)
		if err != nil {
			log.Println("")
			return nil, fmt.Errorf("couldn't parse postgres connection string, %w", err)
		}

		db := pg.Connect(pgOptions)

		_, err = db.Exec("SELECT 1")
		if err != nil {
			return nil, fmt.Errorf("could not connect to your PostgreSQL DB, ERROR: %w", err)
		}

		p := PGDB{db: db}
		err = p.CreateTable()
		if err != nil {
			return nil, fmt.Errorf("could not create migrations table in postgres")
		}

		return p, nil
	case strings.Contains(connString, "mongodb"):
		ctx := context.TODO()
		cs, err := mconnString.ParseAndValidate(connString)
		if err != nil {
			return nil, err
		}
		clientOptions := options.Client().ApplyURI(connString)
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			return nil, fmt.Errorf("could not connect to MongoDB, ERROR: %w", err)
		}
		err = client.Ping(ctx, nil)
		if err != nil {
			return nil, err
		}
		return MongoDb{coll: *client.Database(cs.Database).Collection("migrations")}, nil
	}

	var jsonFile afero.File
	var err error
	jsonFile, err = FS.OpenFile(wd+"/"+dir+"/migrate.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return File{file: jsonFile}, nil
}

func TrackToSetDown(track types.Track, num int) types.Track {
	if len(track.Migrations) == 0 || num <= 0 || num > len(track.Migrations) {
		return types.Track{}
	}

	newLen := len(track.Migrations) - num
	track.Migrations = track.Migrations[:newLen]

	if newLen > 0 {
		last := track.Migrations[newLen-1]
		track.LastRun = last.Title
	} else {
		return types.Track{}
	}

	return track
}

func ProcessLogs(logs string) (types.Track, int) {
	track := types.Track{}
	lines := strings.Split(logs, "\n")
	num := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Only consider up migrations for tracking
		if !strings.HasPrefix(line, "__COMPLETE__[up]:") {
			continue
		}

		filename := strings.TrimSpace(strings.TrimPrefix(line, "__COMPLETE__[up]:"))

		migration := types.Migration{
			Title:      filename,
			ExecutedAt: time.Now().Format("02-01-2006 15:04"),
		}

		track.Migrations = append(track.Migrations, migration)
		track.LastRun = filename
		num++
	}

	return track, num
}
