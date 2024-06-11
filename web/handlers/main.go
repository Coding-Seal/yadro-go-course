package handlers

import (
	"log/slog"
	"net/http"

	"yadro-go-course/web/templates"
)

func MainHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("main handler")

		err := templates.Main(w, templates.MainParams{Layout: templates.Layout{PageTitle: "Main Page"}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			slog.Error("error in main handler", slog.Any("error", err))
		}
	})
}
