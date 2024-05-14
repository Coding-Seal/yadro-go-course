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

func NewFetcher(fetcher ports.ComicFetcherRepo, repo ports.ComicsRepo, search ports.SearchComicsRepo) *Fetcher {
	return &Fetcher{
		repo:    repo,
		search:  search,
		fetcher: fetcher,
	}
}

func (f *Fetcher) Update(ctx context.Context) error { // TODO: revisit
	lastID, err := f.fetcher.LastComicID(ctx)
	if err != nil {
		return err
	}

	var missing []int

	for i := 1; i <= lastID; i++ {
		if _, err := f.repo.Comic(ctx, i); err != nil {
			missing = append(missing, i)
		}
	}

	ids, res := f.fetcher.Comics(ctx, len(missing))

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
