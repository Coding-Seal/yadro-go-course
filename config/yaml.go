package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		DB      `yaml:"db"`
		Fetcher `yaml:"fetcher"`
		Server  `yaml:"server,omitempty"`
		Logger  `yaml:"logger"`
		Search  `yaml:"search"`
	}
	Logger struct {
		Type  string `yaml:"type"`
		Level string `yaml:"level"`
	}
	DB struct {
		Type    string `yaml:"type"`
		Url     string `yaml:"url"`
		Version uint   `yaml:"version"`
	}
	Search struct {
		StopWordsFile string `yaml:"stop_words_file"`
	}
	Fetcher struct {
		SourceURL  string `yaml:"source_url"`
		Parallel   int    `yaml:"parallel"`
		UpdateSpec string `yaml:"update_spec"`
	}
	Server struct {
		Port        int           `yaml:"port"`
		RateLimit   int           `yaml:"rate_limit"`
		DeleteEvery time.Duration `yaml:"delete_every"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	config := &Config{
		DB: DB{
			Type: "json",
			Url:  "db.jsonl",
		},
		Server: Server{
			Port: 8080,
		},
		Logger: Logger{
			Type:  "json",
			Level: "debug",
		},
		Search: Search{
			StopWordsFile: "stopwords.txt",
		},
		Fetcher: Fetcher{
			SourceURL:  "https://xkcd.com",
			Parallel:   100,
			UpdateSpec: "0 0 * * ?",
		},
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
