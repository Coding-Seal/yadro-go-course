package app

import (
	"context"
	"log"
)

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

	for i := 0; i < len(missingComics); i++ {
		fetchedComic, ok := <-comics
		if fetchedComic.Err() == nil && ok {
			c := a.toComic(fetchedComic.Comic)

			a.comics[fetchedComic.Comic.ID] = c

			if err := a.db.Save(c); err != nil {
				log.Println("error appending comic: ", err)
			}
		}
	}

	return nil
}

func (a *App) FetchLastComicID(ctx context.Context) (int, error) {
	return a.fetcher.LastID(ctx)
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
		fetchedComic, ok := <-comics
		if fetchedComic.Err() == nil && ok {
			a.comics[fetchedComic.Comic.ID] = a.toComic(fetchedComic.Comic)
		}
	}

	return nil
}
