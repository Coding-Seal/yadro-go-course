package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	http_util "yadro-go-course/pkg/http-util"

	"yadro-go-course/web/rest"
	"yadro-go-course/web/templates"
)

func Login() http_util.ErrHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := templates.Login(w, templates.LoginParams{Layout: templates.Layout{PageTitle: "Login"}})
		if err != nil {
			return errors.Join(err, http_util.ErrInternal)
		}

		return nil
	}
}

func LoginForm(api *rest.Client) http_util.ErrHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		slog.Debug("LoginForm", slog.String("url", r.URL.Path))
		login := r.FormValue("login")
		pswd := r.FormValue("pswd")

		if pswd == "" || login == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return http_util.ErrNoLoginOrPassword
		}

		token, err := api.Login(r.Context(), login, pswd)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return http_util.ErrForbidden
		}

		cookie := &http.Cookie{
			Name:     "Authorization",
			Value:    token,
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return nil
	}
}
