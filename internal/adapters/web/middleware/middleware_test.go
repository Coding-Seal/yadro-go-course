package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"yadro-go-course/internal/contextutil"
)

func TestChain(t *testing.T) {
	c := Chain(AddRequestID, Logging)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	c(dummyHandler).ServeHTTP(w, r)
}

func TestLogging(t *testing.T) {
	ctx := contextutil.WithReqID(context.Background(), 1)
	h := Logging(dummyHandler)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	h.ServeHTTP(w, r)
}

func TestAddRequestID(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h := AddRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotPanics(t, func() { contextutil.ReqID(r.Context()) })
	}))
	h.ServeHTTP(w, r)
}

func TestWrappedResponseWriter_WriteHeader(t *testing.T) {
	w := wrappedResponseWriter{
		ResponseWriter: httptest.NewRecorder(),
		statusCode:     0,
	}
	w.WriteHeader(http.StatusTeapot)
	assert.Equal(t, http.StatusTeapot, w.statusCode)
}
