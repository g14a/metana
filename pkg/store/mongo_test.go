package store

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/g14a/metana/pkg/types"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mconnString "go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

func TestMongoDb_Set(t *testing.T) {
	mongoUrl := os.Getenv("MONGO_TEST_URL")
	assert.NotEmpty(t, mongoUrl, "MONGO_TEST_URL env var must be set")

	ctx := context.TODO()
	cs, err := mconnString.ParseAndValidate(mongoUrl)
	assert.NoError(t, err)

	clientOptions := options.Client().ApplyURI(mongoUrl)
	client, err := mongo.Connect(ctx, clientOptions)
	assert.NoError(t, err)

	err = client.Ping(ctx, nil)
	assert.NoError(t, err)

	// Clean up test collection
	coll := client.Database(cs.Database).Collection("migrations")
	_ = coll.Drop(ctx)

	m := MongoDb{coll: *coll}

	track := types.Track{
		LastRun: "1745037825_initSchema3.go",
		Migrations: []types.Migration{
			{
				Title:      "1745037800_initSchema.go",
				ExecutedAt: time.Now().Format("02-01-2006 15:04"),
			},
			{
				Title:      "1745037824_initSchema2.go",
				ExecutedAt: time.Now().Format("02-01-2006 15:04"),
			},
			{
				Title:      "1745037825_initSchema3.go",
				ExecutedAt: time.Now().Format("02-01-2006 15:04"),
			},
		},
	}

	err = m.Set(track, afero.NewMemMapFs())
	assert.NoError(t, err)

	actualTrackSet, err := m.Load(afero.NewMemMapFs())
	assert.NoError(t, err)

	assert.Equal(t, track.LastRun, actualTrackSet.LastRun)
	assert.Equal(t, len(track.Migrations), len(actualTrackSet.Migrations))

	for i, m := range track.Migrations {
		assert.Equal(t, m.Title, actualTrackSet.Migrations[i].Title)
		assert.Equal(t, m.ExecutedAt, actualTrackSet.Migrations[i].ExecutedAt)
	}
}
