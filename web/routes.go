package web

import (
	"embed"
	"io/fs"
	"net/http"

	httputil "yadro-go-course/pkg/http-util"

	"yadro-go-course/web/handlers"
	"yadro-go-course/web/rest"
)

//go:embed templates/static

var static embed.FS

func Routes(client *rest.Client) http.Handler {
	static, _ := fs.Sub(static, "templates/static")
	mux := http.NewServeMux()
	st := httputil.Chain(httputil.AddRequestID, httputil.Logging)

	mux.Handle("POST /loginform", st(httputil.WrapHandler(handlers.LoginForm(client))))
	mux.Handle("GET /login", st(httputil.WrapHandler(handlers.Login())))
	mux.Handle("GET /comics", st(httputil.WrapHandler(handlers.Pics(client))))
	mux.Handle("GET /", st(httputil.WrapHandler(handlers.MainHandler())))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))

	return mux
}
