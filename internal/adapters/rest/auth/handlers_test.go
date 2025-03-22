package auth

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	http_util "yadro-go-course/pkg/http-util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"yadro-go-course/internal/core/models"
	"yadro-go-course/internal/core/ports"
	"yadro-go-course/internal/core/services"
)

type userRepoMock struct {
	mock.Mock
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

func TestLogin_Happy(t *testing.T) {
	var m userRepoMock

	password := "bob"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := models.User{
		ID:       5,
		Login:    "bob",
		Password: hashed,
		IsAdmin:  true,
	}
	m.On("UserLogin", "bob").Return(u, nil)
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("{\"login\":\"%s\", \"password\":\"%s\"}", u.Login, password))
	r := httptest.NewRequest("POST", "/", buf)
	w := httptest.NewRecorder()
	h := Login(services.NewUserService(&m), time.Minute)
	err := h(w, r)
	assert.NoError(t, err)
}

func TestLogin_NoSuchUser(t *testing.T) {
	var m userRepoMock

	password := "bob"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := models.User{
		ID:       5,
		Login:    "bob",
		Password: hashed,
		IsAdmin:  true,
	}

	m.On("UserLogin", "bob").Return(models.User{}, ports.ErrNotFound)

	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("{\"login\":\"%s\", \"password\":\"%s\"}", u.Login, password))
	r := httptest.NewRequest("POST", "/", buf)
	w := httptest.NewRecorder()
	h := Login(services.NewUserService(&m), time.Minute)
	err := h(w, r)
	assert.ErrorIs(t, err, http_util.ErrNotFound)
}

func TestLogin_WrongPassword(t *testing.T) {
	var m userRepoMock

	password := "12345"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := models.User{
		ID:       5,
		Login:    "bob",
		Password: hashed,
		IsAdmin:  true,
	}
	m.On("UserLogin", "bob").Return(u, nil)
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("{\"login\":\"%s\", \"password\":\"%s\"}", u.Login, "bob"))
	r := httptest.NewRequest("POST", "/", buf)
	w := httptest.NewRecorder()
	h := Login(services.NewUserService(&m), time.Minute)
	err := h(w, r)
	assert.ErrorIs(t, err, http_util.ErrForbidden)
}

func TestLogin_NoLogin(t *testing.T) {
	var m userRepoMock

	password := "12345"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := models.User{
		ID:       5,
		Login:    "bob",
		Password: hashed,
		IsAdmin:  true,
	}
	m.On("UserLogin", "bob").Return(u, nil)
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("{\"login\":\"%s\", \"password\":\"%s\"}", "", ""))
	r := httptest.NewRequest("POST", "/", buf)
	w := httptest.NewRecorder()
	h := Login(services.NewUserService(&m), time.Minute)
	err := h(w, r)
	assert.ErrorIs(t, err, http_util.ErrBadRequest)
}
