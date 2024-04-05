package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	DBfile    string `yaml:"db_file"`
	SourceURL string `yaml:"source_url"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{
		DBfile:    "database.json",
		SourceURL: "https://xkcd.com",
	}
	file, err := os.Open(configPath)

	if err != nil {
		return config, err
	}

	defer file.Close()
	d := yaml.NewDecoder(file)
	err = d.Decode(&config)

	return config, err
}
