package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"yadro-go-course/pkg/config"
	"yadro-go-course/pkg/database"
	"yadro-go-course/pkg/xkcd"
)

func main() {
	var printTerm bool

	var numComics int

	flag.BoolVar(&printTerm, "o", false, "Print in terminal")
	flag.IntVar(&numComics, "n", 2917, "How many comics to retrieve")
	flag.Parse()

	conf, err := config.NewConfig("config.yaml")
	if err != nil {
		log.Println("Could not parse config.yaml ", err)
	}

	dbFile, err := os.OpenFile(conf.DBfile, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		log.Println("Could not open db file", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	runChan := make(chan os.Signal, 1)
	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	go func() {
		<-runChan
		cancel()
		log.Println("Shutting down pres Ctrl + c to force")
	}()

	client := xkcd.NewClient(conf.SourceURL)

	log.Println("Starting ...")

	comics := client.GetFirst(ctx, numComics)

	log.Println("Writing to disk ...")

	jsonDB := database.NewJsonDB(dbFile)
	jsonDB.Save(comics)

	if printTerm {
		log.Println("Printing ...")

		textDB := database.NewTextDB(os.Stdout)
		textDB.Save(comics)
	}

	log.Println("Done ...")
}
