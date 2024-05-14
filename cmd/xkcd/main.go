package main

import (
	"flag"
	"log"
	"yadro-go-course/config"
	"yadro-go-course/internal/adapters/web"
)

func main() {
	var configPath string

	var port int

	flag.StringVar(&configPath, "c", "config.yaml", "Path to config file")
	flag.IntVar(&port, "p", 0, "Port to listen on")
	flag.Parse()

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	if port != 0 {
		cfg.Port = port
	}

	web.Run(cfg)
}
