package handlers

import (
	"errors"
	"net/http"

	http_util "yadro-go-course/pkg/http-util"

	"yadro-go-course/web/templates"
)

func MainHandler() http_util.ErrHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := templates.Main(w, templates.MainParams{Layout: templates.Layout{PageTitle: "Main Page"}})
		if err != nil {
			return errors.Join(err, http_util.ErrInternal)
		}

		return nil
	}
}
