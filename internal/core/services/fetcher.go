package services

import (
	"context"
	"yadro-go-course/internal/core/ports"
)

type Fetcher struct {
	fetcher ports.ComicFetcherRepo
	repo    ports.ComicsRepo
	search  ports.SearchComicsRepo
}

func (f *Fetcher) Update(ctx context.Context) error { // TODO: revisit
	last, err := f.fetcher.LastComic(ctx)
	if err != nil {
		return err
	}
	var missing []int
	for i := 1; i <= last.ID; i++ {
		comic, err := f.repo.Comic(ctx, i)
		if err != nil {
			return err
		}
		missing = append(missing, comic.ID)
	}
	ids, res := f.fetcher.Comics(ctx, last.ID)
	for _, id := range missing {
		ids <- id
	}
	for i := 0; i < len(missing); i++ {
		comic, ok := <-res
		if !ok {
			break
		}
		if err = f.repo.Store(ctx, comic); err != nil {
			return err
		}
		f.search.AddComic(ctx, comic)
	}
	return nil
}
