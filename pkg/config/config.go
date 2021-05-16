package config

import (
	"log"

	"github.com/spf13/afero"

	"gopkg.in/yaml.v2"
)

type MetanaConfig struct {
	Dir       string `yaml:"dir"`
	StoreConn string `yaml:"store"`
}

func GetMetanaConfig(FS afero.Fs) (*MetanaConfig, error) {
	var MetanaConfigInstance MetanaConfig

	yamlFile, err := afero.ReadFile(FS, ".metana.yml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &MetanaConfigInstance)
	if err != nil {
		return nil, err
	}
	return &MetanaConfigInstance, nil
}

func SetMetanaConfig(mc *MetanaConfig, FS afero.Fs) error {
	b, err := yaml.Marshal(&mc)
	if err != nil {
		log.Fatal(err)
	}

	err = afero.WriteFile(FS, ".metana.yml", b, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
