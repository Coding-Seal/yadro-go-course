package handlers

import (
	"context"
	"math/rand/v2"
	"net/http"
	"yadro-go-course/internal/core/services"
)

func Update(fetcherSrv services.Fetcher) ErrHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.WithValue(r.Context(), "request_id", rand.Int()) // TODO: create middleware
		err := fetcherSrv.Update(ctx)
		if err != nil {
			return nil // FIXME: ErrInternal
		}
		w.WriteHeader(http.StatusOK)
		return err
	}
}
func Search(searchSrv services.Search) ErrHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.WithValue(r.Context(), "request_id", rand.Int()) // TODO: create middleware
		phrase := r.URL.Query().Get("search")
		comics := searchSrv.SearchComics(ctx, phrase, 10)
		if len(comics) == 0 {
			return nil // FIXME: ErrNotFound
		}
		urls := make([]string, len(comics))
		for _, comic := range comics {
			urls = append(urls, comic.ImgURL)
		}
		return writeJson(w, http.StatusOK, urls)
	}
}
