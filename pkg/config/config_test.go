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

	afero.WriteFile(FS, "/Users/g14a/metana/.metana.yml", []byte(config), 0644)

	mc, err := GetMetanaConfig(FS, "/Users/g14a/metana")
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

	err := SetMetanaConfig(&mc, FS, "/Users/g14a/metana")
	if err != nil {
		fmt.Println(err)
	}

	fileBytes, err := afero.ReadFile(FS, "/Users/g14a/metana/.metana.yml")
	if err != nil {
		return
	}

	assert.Equal(t, "dir: migrations\nstore: migrate.json\n", string(fileBytes))
}
