package app

import (
	"fmt"
	"os"
	"yadro-go-course/internal/comic"
	"yadro-go-course/internal/database"
	"yadro-go-course/pkg/words"
	"yadro-go-course/pkg/xkcd"
)

type App struct {
	fetcher *xkcd.Fetcher
	stemmer *words.Stemmer
	db      *database.JsonDB
	comics  map[int]*comic.Comic
	index   map[string][]int
}

func NewApp(sourceURL string, file *os.File, stopWords map[string]struct{}, concurrencyLimit int) *App {
	return &App{
		fetcher: xkcd.NewFetcher(sourceURL, concurrencyLimit),
		stemmer: words.NewStemmer(stopWords),
		db:      database.NewJsonDB(file),
		comics:  make(map[int]*comic.Comic),
	}
}

func (a *App) PrintComics(num int) {
	i := 1

	for id, c := range a.comics {
		fmt.Printf("ID=%d ImgURl=%s Keywords=%v\n", id, c.ImgURL, c.Keywords)

		i++

		if i > num {
			break
		}
	}
}
func (a *App) PrintAllComics() {
	for id, c := range a.comics {
		fmt.Printf("ID=%d ImgURl=%s Keywords=%v\n", id, c.ImgURL, c.Keywords)
	}
}
