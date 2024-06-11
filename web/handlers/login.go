package handlers

import (
	"log/slog"
	"net/http"

	"yadro-go-course/web/rest"
	"yadro-go-course/web/templates"
)

func Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Login", slog.String("url", r.URL.Path))

		err := templates.Login(w, templates.LoginParams{Layout: templates.Layout{PageTitle: "Login"}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			slog.Error("Login", err)
		}
	})
}

func LoginForm(api *rest.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("LoginForm", slog.String("url", r.URL.Path))
		login := r.FormValue("login")
		pswd := r.FormValue("pswd")

		if pswd == "" || login == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		token, err := api.Login(r.Context(), login, pswd)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		cookie := &http.Cookie{
			Name:     "Authorization",
			Value:    token,
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
