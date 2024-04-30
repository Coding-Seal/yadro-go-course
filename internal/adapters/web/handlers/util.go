package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrHandleFunc func(w http.ResponseWriter, r *http.Request) error

func WrapHandler(fn ErrHandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil { // TODO: log
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func writeJson(w http.ResponseWriter, statusCode int, v any) error {
	en := json.NewEncoder(w)
	err := en.Encode(v)
	if err != nil {
		return err
	}
	w.WriteHeader(statusCode)
	return nil
}
