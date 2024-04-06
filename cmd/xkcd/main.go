package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"yadro-go-course/internal/app"
	"yadro-go-course/pkg/config"
)

func main() {
	var printTerm bool

	var numComics int

	flag.BoolVar(&printTerm, "o", false, "Print in terminal")
	flag.IntVar(&numComics, "n", 0, "How many comics to retrieve")
	flag.Parse()

	conf, err := config.NewConfig("config.yaml")
	if err != nil {
		log.Println("Could not parse config.yaml ", err)
	}

	dbFile, err := os.OpenFile(conf.DBfile, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		log.Fatalln("Could not open db file", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client := app.NewApp(conf.SourceURL, dbFile, nil)

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

	log.Println("Done ...")
}
