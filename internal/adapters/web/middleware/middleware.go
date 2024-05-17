package middleware

import (
	"log/slog"
	"math/rand/v2"
	"net/http"
	"time"

	"yadro-go-course/internal/contextutil"
)

type Middleware func(http.Handler) http.Handler

func Stack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}

		return next
	}
}

func AddRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := contextutil.WithReqID(r.Context(), rand.Int()) // Maybe swap for uuid
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &wrappedResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(ww, r)

		end := time.Since(start)

		reqID, err := contextutil.ReqID(r.Context())
		if err != nil {
			slog.Error("middleware: logging", slog.Any("error", err))
		}

		slog.Debug("middleware: logging",
			slog.Int("req_id", reqID),
			slog.String("method", r.Method),
			slog.String("url", r.RequestURI),
			slog.Int("status", ww.statusCode),
			slog.String("duration", end.String()))
	})
}
