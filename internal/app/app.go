package app

import (
	"context"
	"fmt"
	"maps"
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
}

func NewApp(source string, file *os.File, stopWords map[string]struct{}, concurrencyLimit int) *App {
	return &App{
		fetcher: xkcd.NewFetcher(source, concurrencyLimit),
		stemmer: words.NewStemmer(stopWords),
		db:      database.NewJsonDB(file),
		comics:  make(map[int]*comic.Comic),
	}
}

func (a *App) LoadComics() {
	maps.Copy(a.comics, a.db.Read())
}
func (a *App) SaveComics() {
	a.db.Save(a.comics)
}
func (a *App) FetchRemainingComics(ctx context.Context) error {
	var missingComics []int

	lastID, err := a.FetchLastComicID(ctx)

	if err != nil {
		return err
	}

	for id := 1; id <= lastID; id++ {
		if _, ok := a.comics[id]; !ok {
			missingComics = append(missingComics, id)
		}
	}
	ids, comics := a.fetcher.Comics(ctx, len(missingComics))

	for _, id := range missingComics {
		ids <- id
	}
	close(ids)
	for i := 0; i < lastID; i++ {
		fetchedComic := <-comics
		if fetchedComic.Err() == nil {
			a.comics[fetchedComic.Comic.ID] = a.toComic(fetchedComic.Comic)
		}
	}
	return nil
}

func (a *App) FetchLastComicID(ctx context.Context) (int, error) {
	return a.fetcher.GetLastID(ctx)
}
func (a *App) FetchAll(ctx context.Context) error {
	lastID, err := a.FetchLastComicID(ctx)

	if err != nil {
		return err
	}
	ids, comics := a.fetcher.Comics(ctx, lastID)

	for id := 1; id <= lastID; id++ {
		ids <- id
	}
	close(ids)

	for i := 0; i < lastID; i++ {
		fetchedComic := <-comics
		if fetchedComic.Err() == nil {
			a.comics[fetchedComic.Comic.ID] = a.toComic(fetchedComic.Comic)
		}
	}
	return nil
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
func (a *App) toComic(c *xkcd.Comic) *comic.Comic {
	return &comic.Comic{
		ImgURL:   c.ImgURL,
		Keywords: a.stemmer.Stem(words.ParsePhrase(c.Title + " " + c.AltTranscription + " " + c.Transcription + " " + c.News)),
	}
}
