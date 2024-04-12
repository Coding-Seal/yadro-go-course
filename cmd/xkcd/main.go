package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"yadro-go-course/internal/app"
	"yadro-go-course/pkg/config"
	"yadro-go-course/pkg/words"
)

func main() {
	var stopWords string

	var configName string

	flag.StringVar(&configName, "c", "config.yaml", "Path to config file")
	flag.StringVar(&stopWords, "file", "", "Provide list of stop words")
	flag.Parse()

	conf, err := config.NewConfig(configName)
	if err != nil {
		log.Println("Could not parse", configName, err)
	}

	dbFile, err := os.OpenFile(conf.DBfile, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		log.Fatalln("Could not open db file", err)
	}

	defer dbFile.Close()

	var stopWordsMap map[string]struct{}

	if stopWords != "" {
		stopWordsFile, err := os.Open(stopWords)

		if err != nil {
			log.Fatalln("Could not open stop words file", err)
		}

		stopWordsMap = words.ParseStopWords(stopWordsFile)

		defer stopWordsFile.Close()
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client := app.NewApp(conf.SourceURL, dbFile, stopWordsMap, conf.ConcurrencyLimit)

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
