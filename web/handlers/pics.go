package handlers

import (
	"errors"
	"net/http"

	httputil "yadro-go-course/pkg/http-util"

	"yadro-go-course/web/rest"
	"yadro-go-course/web/templates"
)

func Pics(c *rest.Client) httputil.ErrHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		search := r.FormValue("query")

		pics, err := c.SearchPics(r.Context(), search)
		if err != nil {
			return errors.Join(err, httputil.ErrInternal)
		}

		err = templates.Pics(w, templates.PicsParams{
			Urls:   pics,
			Layout: templates.Layout{},
		})
		if err != nil {
			return errors.Join(err, httputil.ErrInternal)
		}

		return nil
	}
}
