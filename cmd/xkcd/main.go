package main

import (
	"flag"
	"log"
	"yadro-go-course/config"
	"yadro-go-course/internal/adapters/web"
)

func main() {
	var configName string

	var searchPhrase string

	var useIndex bool

	var numComics int

	flag.StringVar(&configName, "c", "config.yaml", "Path to config file")
	flag.StringVar(&searchPhrase, "s", "", "Search words")
	flag.BoolVar(&useIndex, "i", false, "Use index")
	flag.IntVar(&numComics, "n", 10, "Number of comics to print")
	flag.Parse()

	cfg, err := config.NewConfig(configName)
	if err != nil {
		log.Fatalln(err)
	}

	web.Run(cfg)
}
