package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"yadro-go-course/internal/contextutil"
)

func TestWrapHandler(t *testing.T) {
	ctx := contextutil.WithReqID(context.Background(), 1)
	h := WrapHandler(func(w http.ResponseWriter, r *http.Request) error {
		return ErrForbidden
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	h.ServeHTTP(w, r)
	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)

	h = WrapHandler(func(w http.ResponseWriter, r *http.Request) error {
		return ErrBadRequest
	})
	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	h.ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	h = WrapHandler(func(w http.ResponseWriter, r *http.Request) error {
		return ErrNotFound
	})
	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	h.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)

	h = WrapHandler(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	})
	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	h.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
