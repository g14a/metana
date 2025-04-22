package store

import (
	"os"
	"testing"
	"time"

	"github.com/g14a/metana/pkg/types"
	"github.com/go-pg/pg/v10"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestPGDB_Set(t *testing.T) {
	pgUrl := os.Getenv("POSTGRES_TEST_URL")
	assert.NotEmpty(t, pgUrl, "POSTGRES_TEST_URL env var must be set")

	pgOptions, err := pg.ParseURL(pgUrl)
	assert.NoError(t, err)

	db := pg.Connect(pgOptions)
	defer db.Close()

	p := PGDB{db: db}

	// Clean up before test
	_ = p.CreateTable()

	now := time.Now().Format("02-01-2006 15:04")

	track := types.Track{
		LastRun: "1745037825_initSchema3.go",
		Migrations: []types.Migration{
			{Title: "1745037800_initSchema.go", ExecutedAt: now},
			{Title: "1745037824_initSchema2.go", ExecutedAt: now},
			{Title: "1745037825_initSchema3.go", ExecutedAt: now},
		},
	}

	err = p.Set(track, afero.NewMemMapFs())
	assert.NoError(t, err)

	actualTrackSet, err := p.Load(afero.NewMemMapFs())
	assert.NoError(t, err)

	assert.Equal(t, track.LastRun, actualTrackSet.LastRun)
	assert.Equal(t, len(track.Migrations), len(actualTrackSet.Migrations))

	for i := range track.Migrations {
		assert.Equal(t, track.Migrations[i].Title, actualTrackSet.Migrations[i].Title)
		assert.Equal(t, track.Migrations[i].ExecutedAt, actualTrackSet.Migrations[i].ExecutedAt)
	}
}
