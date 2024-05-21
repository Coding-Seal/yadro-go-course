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

func Routes(fetcher *services.Fetcher, search *services.Search, user *services.UserService, limit int, deleteEvery time.Duration, ctx context.Context) http.Handler {
	mux := http.NewServeMux()
	st := middleware.Chain(middleware.AddRequestID, middleware.Logging)
	authSt := middleware.Chain(st, auth.Authenticate, middleware.RateLimitOnID(limit, deleteEvery, ctx))
	authzSt := middleware.Chain(authSt, auth.Authorize)

	mux.Handle("GET /pics", st(handlers.WrapHandler(handlers.Search(search))))
	mux.Handle("POST /update", authzSt(handlers.WrapHandler(handlers.Update(fetcher))))
	mux.Handle("POST /login", st(handlers.WrapHandler(auth.Login(user))))

	return mux
}
