package services

import (
	"context"
	"slices"
	"yadro-go-course/internal/core/models"
	"yadro-go-course/internal/core/ports"
)

type Search struct { // TODO: add logging
	repo   ports.SearchComicsRepo
	comics ports.ComicsRepo
}

func (s *Search) SearchComics(ctx context.Context, query string, limit int) []models.Comic {
	scored := s.repo.SearchComics(ctx, query)
	limit = min(limit, len(scored))
	found := make([]int, 0, limit)
	i := 0
	for id, _ := range scored {
		if i >= limit {
			break
		}
		found = append(found, id)
	}
	slices.SortFunc(found, func(a, b int) int { // check if valid
		sca, scb := scored[a], scored[b]
		if sca == scb {
			return b - a
		} else if sca < scb {
			return -1
		} else {
			return 1
		}
	})
	result := make([]models.Comic, 0, limit)
	for _, id := range found {
		comic, err := s.comics.Comic(ctx, id)
		if err != nil {

		}
		result = append(result, comic)
	}
	return result
}
func (s *Search) AddComic(ctx context.Context, comic models.Comic) {
	s.repo.AddComic(ctx, comic)
}
