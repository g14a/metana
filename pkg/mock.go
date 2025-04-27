package pkg

import (
	"github.com/g14a/metana/pkg/types"
	"github.com/spf13/afero"
)

type MockStore struct {
	Data map[string]string
}

func (m *MockStore) Load(FS afero.Fs) (types.Track, error) {
	var migrations []types.Migration
	for title, ts := range m.Data {
		migrations = append(migrations, types.Migration{
			Title:      title,
			ExecutedAt: ts,
		})
	}
	return types.Track{Migrations: migrations}, nil
}

func (m *MockStore) Set(track types.Track, FS afero.Fs) error {
	return nil
}
