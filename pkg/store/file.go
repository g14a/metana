package store

import (
	"encoding/json"

	"github.com/g14a/metana/pkg/types"
	"github.com/spf13/afero"
)

type File struct {
	file afero.File
}

func (f File) Set(track types.Track, FS afero.Fs) error {
	bytes, err := json.MarshalIndent(track, "", "	")
	if err != nil {
		return err
	}

	err = afero.WriteFile(FS, f.file.Name(), bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (f File) Load(FS afero.Fs) (types.Track, error) {
	track, err := afero.ReadFile(FS, f.file.Name())
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

func (f File) Wipe(FS afero.Fs) error {
	err := afero.WriteFile(FS, f.file.Name(), nil, 0644)
	if err != nil {
		return err
	}

	return nil
}
