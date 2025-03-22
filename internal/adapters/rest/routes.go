package rest

import (
	"context"
	"net/http"
	"time"

	"yadro-go-course/internal/adapters/rest/handlers"
	"yadro-go-course/pkg/http-util"

	"yadro-go-course/config"
	"yadro-go-course/pkg/ratelimiter"

	"yadro-go-course/internal/adapters/rest/auth"

	"yadro-go-course/internal/adapters/rest/middleware"

	"yadro-go-course/internal/core/services"
)

func Api(fetcher *services.Fetcher, search *services.Search, user *services.UserService, cfg *config.Config, ctx context.Context) http.Handler {
	base := http_util.Chain(http_util.AddRequestID, http_util.Logging, middleware.ConcurrencyLimiter(cfg.ConcurrencyLimit))

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
	authSt := http_util.Chain(auth.Authenticate, auth.Authorize)

	mux := http.NewServeMux()
	mux.Handle("POST /login", http_util.Chain(base, limitOnIP)(http_util.WrapHandler(auth.Login(user, cfg.TokenMaxTime))))
	mux.Handle("GET /pics", http_util.Chain(base, limitOnIP)(http_util.WrapHandler(handlers.Search(search))))
	mux.Handle("POST /update", http_util.Chain(base, authSt)(http_util.WrapHandler(handlers.Update(fetcher))))

	return mux
}
