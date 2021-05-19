package wipe

import (
	"fmt"
	"os"
	"testing"

	"github.com/g14a/metana/pkg"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestWipe(t *testing.T) {
	FS := afero.NewMemMapFs()

	FS.MkdirAll("/Users/g14a/metana/migrations/scripts", 0755)

	afero.WriteFile(FS, "/Users/g14a/metana/migrations/scripts/1621331831-Abc.go", []byte("{}"), 0644)
	afero.WriteFile(FS, "/Users/g14a/metana/migrations/scripts/1621331861-Random.go", []byte("{}"), 0644)
	afero.WriteFile(FS, "/Users/g14a/metana/migrations/scripts/1621331869-AddIndexes.go", []byte("{}"), 0644)
	afero.WriteFile(FS, "/Users/g14a/metana/migrations/scripts/1621331874-InitSchema.go", []byte("{}"), 0644)

	os.Chdir("/Users/g14a/metana")

	afero.WriteFile(FS, "migrations/migrate.json", []byte(`{
		"LastRun": "1621331861-Random.go",
		"LastRunTS": 1621331861,
		"Migrations": [
			{
				"title": "1621331831-Abc.go",
				"timestamp": 1621331831
			},
			{
				"title": "1621331861-Random.go",
				"timestamp": 1621331861
			}
		]
	}`), 0644)

	err := Wipe("/Users/g14a/metana", "/Users/g14a/metana", "migrations", "", FS)
	assert.NoError(t, err)

	file, err := afero.ReadFile(FS, "/Users/g14a/metana/migrations/main.go")
	if err != nil {
		fmt.Println(err)
	}

	expectedLines := pkg.GetExpectedLinesInit()

	pkg.ExpectLines(t, string(file), expectedLines...)
}
