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
		Server  `yaml:"server"`
		Logger  `yaml:"logger"`
		Search  `yaml:"search"`
	}
	Logger struct {
		Type  string `yaml:"type,omitempty"`
		Level string `yaml:"level"`
	}
	DB struct {
		Url string `yaml:"url"`
	}
	Search struct {
		StopWordsFile string `yaml:"stop_words_file,omitempty"`
	}
	Fetcher struct {
		SourceURL  string `yaml:"source_url,omitempty"`
		Parallel   int    `yaml:"parallel,omitempty"`
		UpdateSpec string `yaml:"update_spec"`
	}
	Server struct {
		Port             int           `yaml:"port,omitempty"`
		RateLimit        int           `yaml:"rate_limit"`
		DeleteEvery      time.Duration `yaml:"delete_every"`
		ConcurrencyLimit int           `yaml:"concurrency_limit"`
		TokenMaxTime     time.Duration `yaml:"token_max_time"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	config := &Config{
		DB: DB{
			Url: "comics.db",
		},
		Server: Server{
			Port: 8080,
		},
		Logger: Logger{
			Type: "json",
		},
		Search: Search{
			StopWordsFile: "stopwords.txt",
		},
		Fetcher: Fetcher{
			SourceURL: "https://xkcd.com",
			Parallel:  100,
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
