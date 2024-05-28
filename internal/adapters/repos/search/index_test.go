package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"yadro-go-course/internal/core/models"
	"yadro-go-course/internal/core/ports"
	"yadro-go-course/pkg/words"
)

func TestIndex_AddComic(t *testing.T) {
	ctx := context.Background()
	testComic := models.Comic{ID: 42, Title: "TestString and something"}
	index := NewIndex(words.NewStemmer(nil))
	index.AddComic(ctx, testComic)
	found := index.SearchComics(ctx, "TestString and something")
	assert.Contains(t, found, testComic.ID)
}

func TestIndex_Build(t *testing.T) {
	m := newComicRepoMock()
	testComics := []models.Comic{{ID: 1, Title: "testStr"}, {ID: 2, Title: "testStr"}, {ID: 3, Title: "testStr"}}
	m.On("ComicsAll").Return(testComics, nil)

	ctx := context.Background()
	index := NewIndex(words.NewStemmer(nil))
	assert.NoError(t, index.Build(ctx, m))

	found := index.SearchComics(ctx, "testStr")
	for _, comic := range testComics {
		assert.Contains(t, found, comic.ID)
	}
}

type comicRepoMock struct {
	mock.Mock
}

var _ ports.ComicsRepo = (*comicRepoMock)(nil)

func newComicRepoMock() *comicRepoMock {
	return &comicRepoMock{}
}

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
