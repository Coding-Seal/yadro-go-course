package handlers

import (
	"log/slog"
	"net/http"

	"yadro-go-course/web/rest"
	"yadro-go-course/web/templates"
)

func Pics(c *rest.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Comics", slog.String("url", r.URL.Path))
		search := r.FormValue("query")
		pics, _ := c.SearchPics(r.Context(), search)

		err := templates.Pics(w, templates.PicsParams{
			Urls:   pics,
			Layout: templates.Layout{},
		})
		if err != nil {
			slog.Error("Pics template error", slog.Any("error", err))
		}
	})
}
