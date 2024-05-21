package middleware

import (
	"context"
	"log/slog"
	"net"
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

func RateLimitOnIP(limit int, deleteAfter time.Duration, ctx context.Context) Middleware {
	return func(next http.Handler) http.Handler {
		limiter := ratelimiter.NewRateLimiterPerUser[string](limit, deleteAfter, ctx)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if limiter.Allow(ip) {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusTooManyRequests)
				slog.Debug("limited request", slog.Int("req_id", contextutil.ReqID(r.Context())), slog.String("ip", ip))
			}
		})
	}
}
