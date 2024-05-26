package services

import (
	"context"
	"strconv"

	"github.com/stretchr/testify/mock"

	"yadro-go-course/internal/core/models"
	"yadro-go-course/internal/core/ports"
)

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

type fetcherRepoMock struct {
	mock.Mock
}

func newFetcherRepoMock() *fetcherRepoMock {
	return &fetcherRepoMock{}
}

var _ ports.ComicFetcherRepo = (*fetcherRepoMock)(nil)

func (m *fetcherRepoMock) LastComicID(ctx context.Context) (int, error) {
	args := m.Called()
	return args.Get(0).(int), args.Error(1)
}

func (m *fetcherRepoMock) Comics(ctx context.Context, limit int) (chan<- int, <-chan models.Comic) {
	args := m.Called(limit)
	ids := args.Get(0).(chan int)
	comics := args.Get(1).(chan models.Comic)

	go func() {
		for id := range ids {
			comics <- models.Comic{ID: id, Title: strconv.Itoa(id)}
		}

		close(comics)
	}()

	return ids, comics
}

type searchRepoMock struct {
	mock.Mock
}

func newSearchRepoMock() *searchRepoMock {
	return &searchRepoMock{}
}

var _ ports.SearchComicsRepo = (*searchRepoMock)(nil)

func (m *searchRepoMock) SearchComics(ctx context.Context, query string) map[int]int {
	args := m.Called(query)
	return args.Get(0).(map[int]int)
}

func (m *searchRepoMock) AddComic(ctx context.Context, comic models.Comic) {
	m.Called(comic)
}

type userRepoMock struct {
	mock.Mock
}

func newUserRepoMock() *userRepoMock {
	return &userRepoMock{}
}

var _ ports.UserRepo = (*userRepoMock)(nil)

func (m *userRepoMock) UserID(ctx context.Context, id int64) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *userRepoMock) UserLogin(ctx context.Context, login string) (models.User, error) {
	args := m.Called(login)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *userRepoMock) RemoveUser(ctx context.Context, id int64) error {
	return m.Called(id).Error(0)
}

func (m *userRepoMock) AddUser(ctx context.Context, user *models.User) error {
	return m.Called(user).Error(0)
}
