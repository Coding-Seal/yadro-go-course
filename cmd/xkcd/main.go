package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"yadro-go-course/internal/app"
	"yadro-go-course/internal/config"
	"yadro-go-course/pkg/words"
)

func main() {

	var configName string

	flag.StringVar(&configName, "c", "config.yaml", "Path to config file")

	flag.Parse()

	cfg, err := config.NewConfig(configName)
	if err != nil {
		log.Println("Could not parse", configName, err)
	}

	dbFile, err := os.OpenFile(cfg.DBfile, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		log.Fatalln("Could not open db file", err)
	}

	defer dbFile.Close()

	var stopWordsMap map[string]struct{}
	if cfg.StopWordsFile != "" {
		stopWordsFile, err := os.Open(cfg.StopWordsFile)

		if err != nil {
			log.Fatalln("Could not open stop words file", err)
		}
		defer stopWordsFile.Close()

		stopWordsMap = words.ParseStopWords(stopWordsFile)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client := app.NewApp(cfg.SourceURL, dbFile, stopWordsMap, cfg.Parallel)

	log.Println("loading comics from db")
	client.LoadComics()

	log.Println("downloading remaining comics")

	err = client.FetchRemainingComics(ctx)

	if err != nil {
		log.Println("remaining comics fetched failed:", err)
	}

	log.Println("saving comics in DB")
	client.SaveComics()

	log.Println("Done")
}
