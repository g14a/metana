package store

import (
	"context"

	"github.com/g14a/metana/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDb struct {
	coll mongo.Collection
	ctx  context.Context
}

func (m MongoDb) Set(track types.Track) error {
	_, err := m.coll.DeleteMany(m.ctx, bson.M{})
	if err != nil {
		return err
	}
	_, err = m.coll.InsertOne(m.ctx, track)
	if err != nil {
		return err
	}

	return nil
}

func (m MongoDb) Load() (types.Track, error) {
	var track types.Track

	err := m.coll.FindOne(m.ctx, bson.M{}).Decode(&track)
	if err == mongo.ErrNoDocuments {
		return types.Track{}, nil
	}

	return track, nil
}
