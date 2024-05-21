package web

import (
	"context"
	"net/http"
	"time"

	"yadro-go-course/internal/adapters/web/auth"

	"yadro-go-course/internal/adapters/web/middleware"

	"yadro-go-course/internal/adapters/web/handlers"
	"yadro-go-course/internal/core/services"
)

func Routes(fetcher *services.Fetcher, search *services.Search, user *services.UserService, concurrencyLimit, rateLimit int, deleteEvery time.Duration, ctx context.Context) http.Handler {
	mux := http.NewServeMux()
	st := middleware.Chain(middleware.AddRequestID, middleware.Logging)
	authSt := middleware.Chain(st, middleware.ConcurrencyLimiter(concurrencyLimit), auth.Authenticate, middleware.RateLimitOnID(rateLimit, deleteEvery, ctx), auth.Authorize)

	mux.Handle("GET /pics", st(handlers.WrapHandler(handlers.Search(search))))
	mux.Handle("POST /update", authSt(handlers.WrapHandler(handlers.Update(fetcher))))
	mux.Handle("POST /login", st(handlers.WrapHandler(auth.Login(user))))

	return mux
}
