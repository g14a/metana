package config

import (
	"os"

	"github.com/spf13/afero"

	"gopkg.in/yaml.v2"
)

func GetMetanaConfig(FS afero.Fs, wd string) (*MetanaConfig, error) {
	var MetanaConfigInstance MetanaConfig

	_, err := FS.Stat(wd + "/.metana.yml")
	if os.IsNotExist(err) {
		return &MetanaConfigInstance, nil
	}

	if err != nil {
		return nil, err
	}
	yamlFile, err := afero.ReadFile(FS, wd+"/.metana.yml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &MetanaConfigInstance)
	if err != nil {
		return nil, err
	}
	return &MetanaConfigInstance, nil
}

func SetMetanaConfig(mc *MetanaConfig, FS afero.Fs, wd string) error {
	b, err := yaml.Marshal(&mc)
	if err != nil {
		return err
	}

	err = afero.WriteFile(FS, wd+"/.metana.yml", b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func SetEnvironmentMetanaConfig(mc *MetanaConfig, env, store string, FS afero.Fs, wd string) error {
	b, err := yaml.Marshal(&mc)
	if err != nil {
		return err
	}

	err = afero.WriteFile(FS, wd+"/.metana.yml", b, 0644)
	if err != nil {
		return err
	}
	return nil
}

type MetanaConfig struct {
	Dir          string        `yaml:"dir"`
	StoreConn    string        `yaml:"store"`
	Environments []Environment `yaml:"environments"`
}

type Environment struct {
	Name  string `yaml:"name"`
	Dir   string `yaml:"dir"`
	Store string `yaml:"store"`
}
