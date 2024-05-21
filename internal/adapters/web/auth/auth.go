package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"yadro-go-course/internal/contextutil"
)

var jwtSecret = []byte("5902dae04b5ee9cafedfacaf7dbcad276b66e647cb0f62fe7ca3cde2e6351258")

type customClaims struct {
	UserID  int64 `json:"user_id"`
	IsAdmin bool  `json:"is_admin"`
	jwt.RegisteredClaims
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Values("Authorization") != nil {
			tokenString := r.Header.Values("Authorization")[0]

			token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return jwtSecret, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				slog.Debug("failed to validate jwt", slog.Int("req_id", contextutil.ReqID(r.Context())), slog.Any("error", err), slog.String("token", tokenString))
			} else if claims, ok := token.Claims.(*customClaims); ok {
				ctx := contextutil.WithUserID(r.Context(), claims.UserID)
				ctx = contextutil.WithIsAdmin(ctx, claims.IsAdmin)
				r = r.WithContext(ctx)
				slog.Debug("authenticated user", slog.Int("req_id", contextutil.ReqID(r.Context())), slog.Int64("user_id", claims.UserID))

				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				slog.Debug("failed to validate jwt claims", slog.Int("req_id", contextutil.ReqID(r.Context())), slog.String("token", tokenString))
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			slog.Debug("no jwt token", slog.Int("req_id", contextutil.ReqID(r.Context())))
		}
	})
}

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if contextutil.IsAdmin(r.Context()) {
			slog.Debug("authorized user", slog.Int("req_id", contextutil.ReqID(r.Context())), slog.Int64("user_id", contextutil.UserID(r.Context())))
			next.ServeHTTP(w, r)
		} else {
			slog.Debug("failed to authorized user", slog.Int("req_id", contextutil.ReqID(r.Context())), slog.Int64("user_id", contextutil.UserID(r.Context())))
		}
	})
}
