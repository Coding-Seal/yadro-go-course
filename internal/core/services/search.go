package services

import (
	"context"
	"slices"

	"yadro-go-course/internal/core/models"
	"yadro-go-course/internal/core/ports"
)

type Search struct {
	repo   ports.SearchComicsRepo
	comics ports.ComicsRepo
}

func NewSearch(repo ports.SearchComicsRepo, comics ports.ComicsRepo) *Search {
	return &Search{repo: repo, comics: comics}
}

func (s *Search) SearchComics(ctx context.Context, query string, limit int) []models.Comic {
	scored := s.repo.SearchComics(ctx, query)
	limit = min(limit, len(scored))
	found := make([]int, 0, len(scored))

	for id := range scored {
		found = append(found, id)
	}

	slices.SortFunc(found, func(a, b int) int {
		sca, scb := scored[a], scored[b]
		if sca == scb {
			return a - b
		}

		return scb - sca
	})

	result := make([]models.Comic, 0, limit)
	i := 0

	for _, id := range found {
		comic, err := s.comics.Comic(ctx, id)
		if err != nil {
			continue
		}

		if i >= limit {
			break
		}

		i++

		result = append(result, comic)
	}

	return result
}

func (s *Search) AddComic(ctx context.Context, comic models.Comic) {
	s.repo.AddComic(ctx, comic)
}
