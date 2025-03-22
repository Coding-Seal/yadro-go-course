package middleware

import (
	"log/slog"
	"net/http"

	"yadro-go-course/pkg/http-util"

	"yadro-go-course/internal/contextutil"
	"yadro-go-course/pkg/ratelimiter"
)

func RateLimitOnID(limiter *ratelimiter.PerUser[int64]) http_util.Middleware {
	return func(next http.Handler) http.Handler {
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

func RateLimitOnIP(limiter *ratelimiter.PerUser[string]) http_util.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if limiter.Allow(ip) {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusTooManyRequests)
				slog.Debug("limited request", slog.Int("req_id", contextutil.ReqID(r.Context())), slog.String("ip", ip))
			}
		})
	}
}

func ConcurrencyLimiter(limit int) http_util.Middleware {
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
