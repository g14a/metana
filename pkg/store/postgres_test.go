package store

import (
	"os"
	"testing"

	"github.com/g14a/metana/pkg/types"
	"github.com/go-pg/pg/v10"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestPGDB_Set(t *testing.T) {
	pgUrl := os.Getenv("POSTGRES_TEST_URL")
	pgOptions, err := pg.ParseURL(pgUrl)

	db := pg.Connect(pgOptions)

	assert.NoError(t, err)
	p := PGDB{
		db: db,
	}

	track := types.Track{
		LastRun:   "1621746410-AddData.go",
		LastRunTS: 1621746410,
		Migrations: []types.Migration{
			{
				Title:     "1621746399-InitSchema.go",
				Timestamp: 1621746399,
			},
			{
				Title:     "1621746406-AddIndexes.go",
				Timestamp: 1621746406,
			},
			{
				Title:     "1621746410-AddData.go",
				Timestamp: 1621746410,
			},
		},
	}

	err = p.Set(track, afero.NewMemMapFs())
	assert.NoError(t, err)

	actualTrackSet, err := p.Load(afero.NewMemMapFs())
	assert.NoError(t, err)

	assert.Equal(t, track, actualTrackSet)
}
