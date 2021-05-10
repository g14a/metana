package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type MetanaConfig struct {
	Dir       string `yaml:"dir"`
	StoreConn string `yaml:"store"`
}

func GetMetanaConfig() (*MetanaConfig, error) {
	var MetanaConfigInstance MetanaConfig

	yamlFile, err := ioutil.ReadFile(".metana.yml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &MetanaConfigInstance)
	if err != nil {
		return nil, err
	}
	return &MetanaConfigInstance, nil
}
