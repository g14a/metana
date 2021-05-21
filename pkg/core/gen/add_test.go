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

	err := CreateInitConfig("migrations", "github.com/g14a/metana", FS)
	assert.Equal(t, true, err == nil)

	filename, err := CreateMigrationFile("/Users/g14a/metana", "migrations", "initSchema", "", FS)
	assert.Equal(t, true, err == nil)

	err = Regen("migrations", "InitSchema", strings.TrimPrefix(filename, "migrations/scripts/"), true, FS)
	assert.Equal(t, true, err == nil)
}
