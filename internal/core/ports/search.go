package ports

import (
	"context"
	"yadro-go-course/internal/core/models"
)

type (
	SearchComicsRepo interface {
		SearchComics(ctx context.Context, query string) map[int]int
		AddComic(ctx context.Context, comic models.Comic)
	}
	/*	SearchComicsService interface {
		SearchComics(ctx context.Context, query string, limit int) ([]int, error)
		AddComic(ctx context.Context, comic models.Comic) error
	}*/
)
