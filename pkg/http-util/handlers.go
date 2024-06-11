package http_util

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"yadro-go-course/internal/contextutil"
)

type ErrHandleFunc func(w http.ResponseWriter, r *http.Request) error

func WrapHandler(fn ErrHandleFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			if errors.Is(err, ErrInternal) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else if errors.Is(err, ErrNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else if errors.Is(err, ErrBadRequest) {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else if errors.Is(err, ErrForbidden) {
				http.Error(w, err.Error(), http.StatusForbidden)
			}

			slog.Error("Error in handler: ", slog.Int("req_id", contextutil.ReqID(r.Context())), slog.Any("error", err))
		}
	})
}

func WriteJson(w http.ResponseWriter, v any) error {
	en := json.NewEncoder(w)

	err := en.Encode(v)
	if err != nil {
		return err
	}

	return nil
}
