package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

type MetanaConfig struct {
	Dir       string `yaml:"dir"`
	StoreConn string `yaml:"store"`
}

// GetMetanaConfig loads .metana.yml from wd or its parent
func GetMetanaConfig(FS afero.Fs, wd string) (*MetanaConfig, error) {
	var config MetanaConfig

	configPath := filepath.Join(wd, ".metana.yml")
	_, err := FS.Stat(configPath)
	if os.IsNotExist(err) {
		return &config, nil
	}
	if err != nil {
		return nil, err
	}

	data, err := afero.ReadFile(FS, configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func SetMetanaConfig(mc *MetanaConfig, FS afero.Fs, migrationsDir string) error {
	data, err := yaml.Marshal(mc)
	if err != nil {
		return err
	}

	// Assume migrationsDir is migrations or schema-mig
	parentDir := filepath.Dir(migrationsDir)

	configPath := filepath.Join(parentDir, ".metana.yml")

	err = afero.WriteFile(FS, configPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
