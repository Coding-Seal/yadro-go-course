package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"yadro-go-course/internal/contextutil"
	"yadro-go-course/pkg/ratelimiter"
)

func RateLimitOnID(limit int, deleteAfter time.Duration, ctx context.Context) Middleware {
	return func(next http.Handler) http.Handler {
		limiter := ratelimiter.NewRateLimiterPerUser[int64](limit, deleteAfter, ctx)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := contextutil.UserID(r.Context())
			if limiter.Allow(id) {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusTooManyRequests)
				slog.Debug("limited request", slog.Int("req_id", contextutil.ReqID(r.Context())), slog.Int64("user_id", id))
			}
		})
	}
}

func ConcurrencyLimiter(limit int) Middleware {
	return func(next http.Handler) http.Handler {
		sem := make(chan struct{}, limit)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(sem) < cap(sem) {
				sem <- struct{}{}

				next.ServeHTTP(w, r)
				<-sem
			} else {
				w.WriteHeader(http.StatusTooManyRequests)
				slog.Debug("limited request", slog.Int("req_id", contextutil.ReqID(r.Context())))
			}
		})
	}
}
