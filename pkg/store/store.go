package store

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/g14a/metana/pkg/types"
	"github.com/go-pg/pg/v10"
	"github.com/spf13/afero"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mconnString "go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Store interface {
	Set(track types.Track, FS afero.Fs) error
	Load(FS afero.Fs) (types.Track, error)
}

func GetStoreViaConn(connStr, dir string, fs afero.Fs, wd string) (Store, error) {
	connStr = resolveEnvVar(connStr)

	switch {
	case strings.Contains(connStr, "postgres://"):
		return setupPostgres(connStr)
	case strings.Contains(connStr, "mongodb"):
		return setupMongo(connStr)
	default:
		return setupFileStore(fs, filepath(wd, dir, "migrate.json"))
	}
}

func resolveEnvVar(conn string) string {
	if strings.HasPrefix(conn, "@") {
		return os.Getenv(strings.TrimPrefix(conn, "@"))
	}
	return conn
}

func setupPostgres(connStr string) (Store, error) {
	opts, err := pg.ParseURL(connStr)
	if err != nil {
		return nil, fmt.Errorf("invalid postgres URL: %w", err)
	}
	db := pg.Connect(opts)
	if _, err := db.Exec("SELECT 1"); err != nil {
		return nil, fmt.Errorf("postgres connection failed: %w", err)
	}
	p := PGDB{db: db}
	if err := p.CreateTable(); err != nil {
		return nil, fmt.Errorf("failed to create migrations table: %w", err)
	}
	return p, nil
}

func setupMongo(connStr string) (Store, error) {
	cs, err := mconnString.ParseAndValidate(connStr)
	if err != nil {
		return nil, err
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, fmt.Errorf("mongo connection failed: %w", err)
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}
	return MongoDb{coll: *client.Database(cs.Database).Collection("migrations")}, nil
}

func setupFileStore(fs afero.Fs, path string) (Store, error) {
	file, err := fs.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return File{file: file}, nil
}

func filepath(parts ...string) string {
	return strings.Join(parts, "/")
}

func TrackToSetDown(track types.Track, num int) types.Track {
	if len(track.Migrations) < num || num <= 0 {
		return types.Track{}
	}
	track.Migrations = track.Migrations[:len(track.Migrations)-num]
	if len(track.Migrations) == 0 {
		return types.Track{}
	}
	track.LastRun = track.Migrations[len(track.Migrations)-1].Title
	return track
}

func ProcessLogs(logs string) (types.Track, int) {
	track := types.Track{}
	count := 0

	for _, line := range strings.Split(logs, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "__COMPLETE__[up]:") {
			continue
		}

		filename := strings.TrimPrefix(line, "__COMPLETE__[up]:")
		m := types.Migration{
			Title:      strings.TrimSpace(filename),
			ExecutedAt: time.Now().Format("02-01-2006 15:04"),
		}
		track.Migrations = append(track.Migrations, m)
		track.LastRun = m.Title
		count++
	}
	return track, count
}
