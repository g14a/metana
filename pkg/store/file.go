package store

import (
	"encoding/json"
	"fmt"
	"github.com/g14a/go-migrate/pkg/types"
	"os"
)

type File struct {
	file os.File
}

func (f File) Set(track types.Track) error {
	bytes, err := json.MarshalIndent(track, "", "	")
	if err != nil {
		return err
	}

	err = os.WriteFile(f.file.Name(), bytes, 0644)
	if err != nil {
		fmt.Println(err,"===============err")
		return err
	}
	return nil
}

func (f File) Load() (types.Track, error) {
	track, err := os.ReadFile(f.file.Name())
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