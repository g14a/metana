package config

import (
	"fmt"
	"log"

	"github.com/spf13/afero"

	"gopkg.in/yaml.v2"
)

type MetanaConfig struct {
	Dir       string `yaml:"dir"`
	StoreConn string `yaml:"store"`
}

func GetMetanaConfig(FS afero.Fs, wd string) (*MetanaConfig, error) {
	var MetanaConfigInstance MetanaConfig

	fmt.Println(wd + "/.metana.yml")
	yamlFile, err := afero.ReadFile(FS, wd+"/.metana.yml")

	fmt.Println(string(yamlFile),"======string yamlfile=========")

	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &MetanaConfigInstance)
	if err != nil {
		return nil, err
	}
	fmt.Println(&MetanaConfigInstance,"============config========")
	return &MetanaConfigInstance, nil
}

func SetMetanaConfig(mc *MetanaConfig, FS afero.Fs, wd string) error {
	b, err := yaml.Marshal(&mc)
	if err != nil {
		log.Fatal(err)
	}

	err = afero.WriteFile(FS, wd+"/.metana.yml", b, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
