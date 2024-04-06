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
	var printTerm bool

	var numComics int

	var stopWords string

	flag.BoolVar(&printTerm, "o", false, "Print to terminal")
	flag.IntVar(&numComics, "n", 0, "How many comics to print")
	flag.StringVar(&stopWords, "file", "", "Provide list of stop words")
	flag.Parse()

	conf, err := config.NewConfig("config.yaml")
	if err != nil {
		log.Println("Could not parse config.yaml ", err)
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

	client := app.NewApp(conf.SourceURL, dbFile, stopWordsMap)

	client.LoadComics()
	lastID, err := client.FetchLastComicID(ctx)

	if err != nil {
		log.Fatalln("Could not fetch last comic", err)
	}

	client.FetchRemainingComics(lastID, ctx)
	client.SaveComics()

	if printTerm {
		if numComics == 0 {
			client.PrintAllComics()
		} else {
			client.PrintComics(numComics)
		}
	}
}
