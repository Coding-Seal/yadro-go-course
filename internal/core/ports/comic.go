package ports

import (
	"context"

	"yadro-go-course/internal/core/models"
)

type (
	/*	ComicService interface {
		Comic(ctx context.Context, id int) (models.Comic, error)
		Store(ctx context.Context, comic models.Comic) error
	}*/
	ComicsRepo interface {
		Comic(ctx context.Context, id int) (models.Comic, error)
		Store(ctx context.Context, comic models.Comic) error
		ComicsAll(ctx context.Context) ([]models.Comic, error)
	}
)
