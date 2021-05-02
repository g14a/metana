package main

import (
	"encoding/json"
	"os"

	"github.com/g14a/go-migrate/pkg/types"
)

func Set(timestamp int, fileName string, up bool) error {
	track, err := Load()
	if err != nil {
		return err
	}

	if up {
		track.LastRun = fileName
		track.LastRunTS = timestamp
		track.Migrations = append(track.Migrations, types.Migration{
			Title:     fileName,
			Timestamp: timestamp,
		})
	} else {
		if len(track.Migrations) == 0 {
			err = os.WriteFile("migrate.json", nil, 0644)
			if err != nil {
				return err
			}
			return nil
		}
		track.LastRun = fileName
		track.LastRunTS = timestamp
		track.Migrations = track.Migrations[:len(track.Migrations)-1]
	}

	bytes, err := json.MarshalIndent(track, "", "	")
	if err != nil {
		return err
	}

	err = os.WriteFile("migrate.json", bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func Load() (types.Track, error) {
	track, err := os.ReadFile("migrate.json")
	if err != nil {
		return types.Track{}, err
	}

	t := types.Track{}

	if len(track) > 0 {
		err = json.Unmarshal(track, &t)
		if err != nil {
			return types.Track{}, err
		}
	}

	return t, nil
}
