package config

import (
	"fmt"
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
		return config, fmt.Errorf("loading config: %w", err)
	}

	defer file.Close()
	d := yaml.NewDecoder(file)
	err = d.Decode(&config)

	return config, fmt.Errorf("parsing config: %w", err)
}
