package ports

import (
	"context"
	"yadro-go-course/internal/core/models"
)

type (
	ComicFetcherRepo interface {
		LastComicID(ctx context.Context) (int, error)
		Comics(ctx context.Context, limit int) (chan<- int, <-chan models.Comic)
	}
	/*	ComicFetcherService interface {
		Update(ctx context.Context) error
	}*/
)
