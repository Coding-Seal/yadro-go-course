package handlers

import (
	"context"

	"github.com/stretchr/testify/mock"
	"yadro-go-course/internal/core/models"
	"yadro-go-course/internal/core/ports"
)

type comicRepoMock struct {
	mock.Mock
}

var _ ports.ComicsRepo = (*comicRepoMock)(nil)

func (m *comicRepoMock) Comic(ctx context.Context, id int) (models.Comic, error) {
	args := m.Called(id)
	return args.Get(0).(models.Comic), args.Error(1)
}

func (m *comicRepoMock) Store(ctx context.Context, comic models.Comic) error {
	args := m.Called(comic)
	return args.Error(0)
}

func (m *comicRepoMock) ComicsAll(ctx context.Context) ([]models.Comic, error) {
	args := m.Called()
	return args.Get(0).([]models.Comic), args.Error(1)
}

type searchRepoMock struct {
	mock.Mock
}

var _ ports.SearchComicsRepo = (*searchRepoMock)(nil)

func (m *searchRepoMock) SearchComics(ctx context.Context, query string) map[int]int {
	args := m.Called(query)
	return args.Get(0).(map[int]int)
}

func (m *searchRepoMock) AddComic(ctx context.Context, comic models.Comic) {
	m.Called(comic)
}
