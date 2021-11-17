package gen

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestRegen(t *testing.T) {
	FS := afero.NewMemMapFs()

	os.Chdir("/Users/g14a/metana")

	err := CreateInitConfig("migrations", "github.com/g14a/metana", FS, "")
	assert.Equal(t, true, err == nil)

	cOpts := CreateMigrationOpts{
		Wd:            "/Users/g14a/metana",
		MigrationsDir: "migrations",
		File:          "initSchema",
		CustomTmpl:    "",
		Environment:   "",
		FS:            FS,
	}
	filename, err := CreateMigrationFile(cOpts)
	assert.Equal(t, true, err == nil)

	rOpts := RegenOpts{
		MigrationsDir:  "migrations",
		MigrationName:  "InitSchema",
		Filename:       strings.TrimPrefix(filename, "migrations/scripts/"),
		FirstMigration: true,
		Environment:    "",
		FS:             FS,
		GoModPath:      "/Users/g14a/metana",
	}

	err = Regen(rOpts)
	assert.Equal(t, true, err == nil)
}
