package web

import (
	"embed"
	"io/fs"
	"net/http"

	"yadro-go-course/web/handlers"
	"yadro-go-course/web/rest"
)

//go:embed templates/static

var static embed.FS

func Routes(client *rest.Client) http.Handler {
	static, _ := fs.Sub(static, "templates/static")
	mux := http.NewServeMux()
	mux.Handle("POST /loginform", handlers.LoginForm(client))
	mux.Handle("GET /login", handlers.Login())
	mux.Handle("GET /comics", handlers.Pics(client))
	mux.Handle("GET /", handlers.MainHandler())
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))
	return mux
}
