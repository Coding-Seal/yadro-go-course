package web

import (
	"net/http"

	"yadro-go-course/internal/adapters/web/handlers"
	"yadro-go-course/internal/adapters/web/middleware"
	"yadro-go-course/internal/core/services"
)

func Routes(fetcher *services.Fetcher, search *services.Search) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /pics", handlers.WrapHandler(handlers.Search(search)))
	mux.Handle("POST /update", handlers.WrapHandler(handlers.Update(fetcher)))

	st := middleware.Stack(middleware.AddRequestID, middleware.Logging)

	return st(mux)
}
