package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"slices"
	"yadro-go-course/internal/app"
	"yadro-go-course/internal/comic"
	"yadro-go-course/internal/config"
	"yadro-go-course/pkg/words"
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

	if searchPhrase == "" {
		log.Fatalln("No search phrase provided")
	}

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

	if err := client.LoadComics(); err != nil {
		log.Println("Could not load comics:", err)
	}

	log.Println("downloading remaining comics")

	err = client.FetchRemainingComics(ctx)

	if err != nil {
		log.Println("remaining comics fetched failed:", err)
	}

	log.Println("saving comics in DB")
	//client.SaveComics()

	var foundComics map[*comic.Comic]int

	if useIndex {
		client.BuildIndex()
		foundComics = client.SearchIndex(searchPhrase)
	} else {
		foundComics = client.SearchComics(searchPhrase)
	}

	type pair struct {
		comic *comic.Comic
		score int
	}

	foundComicsSorted := make([]pair, 0, len(foundComics))

	for c, score := range foundComics {
		foundComicsSorted = append(foundComicsSorted, pair{comic: c, score: score})
	}

	slices.SortFunc(foundComicsSorted, func(a, b pair) int {
		if a.score == b.score {
			if a.comic.Title < b.comic.Title {
				return -1
			} else if a.comic.Title > b.comic.Title {
				return 1
			} else {
				return 0
			}
		}

		return b.score - a.score
	})

	for i := 0; i < numComics; i++ {
		log.Println(foundComicsSorted[i].comic.Title, foundComicsSorted[i].comic.ImgURL,
			"score:", foundComicsSorted[i].score)
	}

	log.Println("Done")
}
