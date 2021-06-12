package pkg

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGetMigrations_validFiles(t *testing.T) {
	FS := afero.NewMemMapFs()

	FS.MkdirAll("/Users/g14a/metana/migrations/scripts", 0755)
	afero.WriteFile(FS, "/Users/g14a/metana/migrations/scripts/1621081055-InitSchema.go", []byte("{}"), 0644)
	afero.WriteFile(FS, "/Users/g14a/metana/migrations/scripts/1621084125-AddIndexes.go", []byte("{}"), 0644)
	afero.WriteFile(FS, "/Users/g14a/metana/migrations/scripts/1621084135-AddFKeys.go", []byte("{}"), 0644)

	os.Chdir("/Users/g14a/metana")
	migrations, err := GetMigrations("/Users/g14a/metana", "migrations", FS, "")
	if err != nil {
		return
	}
	wantedMigrations := []Migration{
		{
			Name: "1621081055-InitSchema.go",
		},
		{
			Name: "1621084125-AddIndexes.go",
		},
		{
			Name: "1621084135-AddFKeys.go",
		},
	}
	for i, m := range migrations {
		assert.Equal(t, wantedMigrations[i].Name, m.Name)
	}
}

func TestGetMigrations_no_files(t *testing.T) {
	FS := afero.NewMemMapFs()

	FS.MkdirAll("/Users/g14a/metana/migrations/scripts", 0755)

	os.Chdir("/Users/g14a/metana")
	migrations, err := GetMigrations("/Users/g14a/metana", "migrations", FS, "")
	if err != nil {
		return
	}

	assert.Equal(t, 0, len(migrations))
}
