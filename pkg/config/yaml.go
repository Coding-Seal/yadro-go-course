package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	SourceURL string `yaml:"source_url"`
	DBfile    string `yaml:"db_file"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{
		SourceURL: "https://xkcd.com",
		DBfile:    "database.json",
	}
	file, err := os.Open(configPath)

	if err != nil {
		return config, fmt.Errorf("loading config: %w", err)
	}

	defer file.Close()
	d := yaml.NewDecoder(file)
	err = d.Decode(&config)
	if err != nil {
		return config, fmt.Errorf("parsing config: %w", err)
	}

	return config, nil
}
