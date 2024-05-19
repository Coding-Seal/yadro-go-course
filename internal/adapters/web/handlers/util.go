package handlers

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
			}

			reqID, idErr := contextutil.ReqID(r.Context())

			if idErr != nil {
				slog.Error("No req_id in context", slog.String("url", r.URL.String()))
			}

			slog.Error("Error in handler: ", slog.Int("req_id", reqID), slog.Any("error", err))
		}
	})
}

func writeJson(w http.ResponseWriter, v any) error {
	en := json.NewEncoder(w)

	err := en.Encode(v)
	if err != nil {
		return err
	}

	return nil
}
