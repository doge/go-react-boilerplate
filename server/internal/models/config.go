package models

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
		IP       string `yaml:"ip"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
}

func (c Config) URI() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		c.Database.Username, c.Database.Password,
		c.Database.IP, c.Database.Port,
	)
}

func LoadConfig(configName string) *Config {
	var cfg *Config

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(cwd + "/config/" + configName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}
