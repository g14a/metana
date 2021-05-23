package store

import (
	"context"
	"os"
	"testing"

	"github.com/g14a/metana/pkg/types"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mconnString "go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

func TestMongoDb_Set(t *testing.T) {
	mongoUrl := os.Getenv("MONGO_TEST_URL")
	ctx := context.TODO()
	cs, err := mconnString.ParseAndValidate(mongoUrl)
	assert.NoError(t, err)
	clientOptions := options.Client().ApplyURI(mongoUrl)
	client, err := mongo.Connect(ctx, clientOptions)
	assert.NoError(t, err)
	err = client.Ping(ctx, nil)
	m := MongoDb{
		coll: *client.Database(cs.Database).Collection("migrations"),
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

	err = m.Set(track, afero.NewMemMapFs())
	assert.NoError(t, err)

	actualTrackSet, err := m.Load(afero.NewMemMapFs())
	assert.NoError(t, err)

	assert.Equal(t, track, actualTrackSet)
}
