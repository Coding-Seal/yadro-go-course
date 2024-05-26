package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"yadro-go-course/internal/core/models"
	"yadro-go-course/internal/core/ports"
)

func TestFetcher_Update(t *testing.T) {
	fetcherRepo := newFetcherRepoMock()
	comicRepo := newComicRepoMock()
	searchRepo := newSearchRepoMock()

	fetcherRepo.On("LastComicID").Return(5, nil)
	fetcherRepo.On("Comics", 3).Return(make(chan int, 3), make(chan models.Comic, 3))
	comicRepo.On("Comic", 1).Return(models.Comic{ID: 1, Title: "1"}, nil)
	comicRepo.On("Comic", 2).Return(models.Comic{ID: 2, Title: "2"}, nil)
	comicRepo.On("Comic", 3).Return(models.Comic{}, ports.ErrNotFound)
	comicRepo.On("Comic", 4).Return(models.Comic{}, ports.ErrNotFound)
	comicRepo.On("Comic", 5).Return(models.Comic{}, ports.ErrNotFound)
	comicRepo.On("Store", models.Comic{ID: 3, Title: "3"}).Return(nil)
	comicRepo.On("Store", models.Comic{ID: 4, Title: "4"}).Return(nil)
	comicRepo.On("Store", models.Comic{ID: 5, Title: "5"}).Return(nil)
	searchRepo.On("AddComic", mock.Anything)

	fetcher := NewFetcher(fetcherRepo, comicRepo, searchRepo)

	err := fetcher.Update(context.Background())
	assert.NoError(t, err)
}
