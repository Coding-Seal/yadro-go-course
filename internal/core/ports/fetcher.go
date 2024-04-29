package ports

import (
	"context"
	"yadro-go-course/internal/core/models"
)

type (
	ComicFetcherRepo interface {
		LastComic(ctx context.Context) (models.Comic, error)
		Comics(ctx context.Context, limit int) (chan<- int, <-chan models.Comic)
	}
	ComicFetcherService interface {
		FetchRemainingComics(ctx context.Context) ([]models.Comic, error)
		Update(ctx context.Context) error
	}
)
