package web

import (
	"context"
	"net/http"
	"time"

	"yadro-go-course/config"
	"yadro-go-course/pkg/ratelimiter"

	"yadro-go-course/internal/adapters/web/auth"

	"yadro-go-course/internal/adapters/web/middleware"

	"yadro-go-course/internal/adapters/web/handlers"
	"yadro-go-course/internal/core/services"
)

func Api(fetcher *services.Fetcher, search *services.Search, user *services.UserService, cfg *config.Config, ctx context.Context) http.Handler {
	base := middleware.Chain(middleware.AddRequestID, middleware.Logging, middleware.ConcurrencyLimiter(cfg.ConcurrencyLimit))

	limiter := ratelimiter.NewRateLimiterPerUser[string](cfg.RateLimit, cfg.DeleteEvery)
	ticker := time.NewTicker(cfg.DeleteEvery)

	go func() {
		for {
			select {
			case <-ticker.C:
				limiter.CleanUp()
			case <-ctx.Done():
				return
			}
		}
	}()

	limitOnIP := middleware.RateLimitOnIP(limiter)
	authSt := middleware.Chain(auth.Authenticate, auth.Authorize)

	mux := http.NewServeMux()
	mux.Handle("POST /login", middleware.Chain(base, limitOnIP)(handlers.WrapHandler(auth.Login(user, cfg.TokenMaxTime))))
	mux.Handle("GET /pics", middleware.Chain(base, limitOnIP)(handlers.WrapHandler(handlers.Search(search))))
	mux.Handle("POST /update", middleware.Chain(base, authSt)(handlers.WrapHandler(handlers.Update(fetcher))))

	return mux
}
