package app

import (
	"context"
	"fmt"
	"io"
	"maps"
	"yadro-go-course/pkg/comic"
	"yadro-go-course/pkg/database"
	"yadro-go-course/pkg/words"
	"yadro-go-course/pkg/xkcd"
)

type App struct {
	fetcher *xkcd.Fetcher
	stemmer *words.Stemmer
	db      *database.JsonDB
	comics  map[int]*comic.Comic
}

func NewApp(source string, rw io.ReadWriter, stopWords map[string]struct{}) *App {
	return &App{
		fetcher: xkcd.NewFetcher(source),
		stemmer: words.NewStemmer(stopWords),
		db:      database.NewJsonDB(rw),
		comics:  make(map[int]*comic.Comic),
	}
}

func (a *App) LoadComics() {
	maps.Copy(a.comics, a.db.Read())
}
func (a *App) SaveComics() {
	a.db.Save(a.comics)
}
func (a *App) FetchRemainingComics(lastID int, ctx context.Context) {
	var missingComics []int

	for id := 1; id <= lastID; id++ {
		if _, ok := a.comics[id]; !ok {
			missingComics = append(missingComics, id)
		}
	}

	for _, fetchedComic := range a.fetcher.GetComics(ctx, missingComics) {
		if fetchedComic != nil {
			a.comics[fetchedComic.ID] = a.toComic(fetchedComic)
		}
	}
}

func (a *App) FetchLastComicID(ctx context.Context) (int, error) {
	return a.fetcher.GetLastID(ctx)
}
func (a *App) FetchAll(lastID int, ctx context.Context) {
	for _, fetchedComic := range a.fetcher.GetALLComics(ctx, lastID) {
		if fetchedComic != nil {
			a.comics[fetchedComic.ID] = a.toComic(fetchedComic)
		}
	}
}

func (a *App) PrintComics(num int) {
	i := 1

	for _, c := range a.comics {
		fmt.Printf("ID=%d ImgURl=%s Keywords=%v\n", c.ID, c.ImgURL, c.Keywords)

		i++

		if i >= num {
			break
		}
	}
}
func (a *App) PrintAllComics() {
	for _, c := range a.comics {
		fmt.Printf("ID=%d ImgURl=%s Keywords=%v\n", c.ID, c.ImgURL, c.Keywords)
	}
}
func (a *App) toComic(f *xkcd.FetchedComic) *comic.Comic {
	if f == nil {
		return nil
	}

	return &comic.Comic{
		ID:       f.ID,
		ImgURL:   f.ImgURL,
		Keywords: a.stemmer.Stem(words.ParsePhrase(f.Title + " " + f.AltTranscription + " " + f.Transcription)),
	}
}
