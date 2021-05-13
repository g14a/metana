package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type MetanaConfig struct {
	Dir       string `yaml:"dir"`
	StoreConn string `yaml:"store"`
	Wipe      bool   `yaml:"wipe"`
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

func SetMetanaConfig(mc *MetanaConfig) error {
	b, err := yaml.Marshal(mc)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(".metana.yml", b, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
