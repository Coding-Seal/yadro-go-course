package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"yadro-go-course/internal/contextutil"
)

func TestAuthenticate_NoToken(t *testing.T) {
	ctx := contextutil.WithReqID(context.Background(), 1)
	auth := Authenticate(http.HandlerFunc(dummyHandler))
	r := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	w := httptest.NewRecorder()
	auth.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthenticate_Happy(t *testing.T) {
	ctx := contextutil.WithReqID(context.Background(), 1)
	auth := Authenticate(http.HandlerFunc(dummyHandler))
	r := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	claims := customClaims{UserID: 1, IsAdmin: false, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute))}}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(jwtSecret)
	assert.NoError(t, err)
	r.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	auth.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthenticate_NoExpiration(t *testing.T) {
	ctx := contextutil.WithReqID(context.Background(), 1)
	auth := Authenticate(http.HandlerFunc(dummyHandler))
	r := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	claims := customClaims{UserID: 1, IsAdmin: false}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(jwtSecret)
	assert.NoError(t, err)
	r.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	auth.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestAuthorize_Happy(t *testing.T) {
	ctx := contextutil.WithIsAdmin(contextutil.WithReqID(context.Background(), 1), true)
	ctx = contextutil.WithUserID(ctx, 1)
	auth := Authorize(http.HandlerFunc(dummyHandler))
	r := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	w := httptest.NewRecorder()
	auth.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthorize(t *testing.T) {
	ctx := contextutil.WithIsAdmin(contextutil.WithReqID(context.Background(), 1), false)
	ctx = contextutil.WithUserID(ctx, 1)
	auth := Authorize(http.HandlerFunc(dummyHandler))
	r := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	w := httptest.NewRecorder()
	auth.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
