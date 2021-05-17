package config

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGetMetanaConfig(t *testing.T) {
	FS := afero.NewMemMapFs()

	config := "dir: migrations\nstore: \"\""

	afero.WriteFile(FS, ".metana.yml", []byte(config), 0644)

	mc, err := GetMetanaConfig(FS)
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, mc != nil)
	assert.Equal(t, "migrations", mc.Dir)
	assert.Equal(t, "", mc.StoreConn)
}

func TestSetMetanaConfig(t *testing.T) {
	FS := afero.NewMemMapFs()

	mc := MetanaConfig{
		Dir:       "migrations",
		StoreConn: "migrate.json",
	}

	err := SetMetanaConfig(&mc, FS)
	if err != nil {
		fmt.Println(err)
	}

	fileBytes, err := afero.ReadFile(FS, ".metana.yml")
	if err != nil {
		return
	}

	assert.Equal(t, "dir: migrations\nstore: migrate.json\n", string(fileBytes))
}
