package services

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"yadro-go-course/internal/core/models"
)

func TestSearch_AddComic(t *testing.T) {
	searchRepo := newSearchRepoMock()
	comicRepo := newComicRepoMock()
	search := NewSearch(searchRepo, comicRepo)
	searchRepo.On("AddComic", mock.Anything).Return(nil)
	search.AddComic(context.Background(), models.Comic{ID: 1, Title: "1"})
}

func TestSearch_SearchComics(t *testing.T) {
	searchRepo := newSearchRepoMock()
	comicRepo := newComicRepoMock()
	search := NewSearch(searchRepo, comicRepo)

	searchRepo.On("SearchComics", "1").Return(map[int]int{1: 5, 2: 5, 3: 5, 4: 5, 5: 4, 6: 4, 7: 4, 8: 4})

	expected := make([]models.Comic, 0, 6)

	for i := 1; i <= 8; i++ {
		comicRepo.On("Comic", i).Return(models.Comic{ID: i, Title: strconv.Itoa(i)}, nil)
	}

	for i := 1; i <= 6; i++ {
		expected = append(expected, models.Comic{ID: i, Title: strconv.Itoa(i)})
	}

	comics := search.SearchComics(context.Background(), "1", 6)
	assert.Equal(t, expected, comics)
}
