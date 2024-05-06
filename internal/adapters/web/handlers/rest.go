package handlers

import (
	"errors"
	"net/http"
	"yadro-go-course/internal/core/services"
)

var (
	ErrNotFound = errors.New("not found")
	ErrInternal = errors.New("internal server error")
)

func Update(fetcherSrv *services.Fetcher) ErrHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := fetcherSrv.Update(r.Context())
		if err != nil {
			return errors.Join(ErrInternal, err)
		}

		return nil
	}
}
func Search(searchSrv *services.Search) ErrHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		phrase := r.URL.Query().Get("search")
		comics := searchSrv.SearchComics(r.Context(), phrase, 10)

		if len(comics) == 0 {
			return ErrNotFound
		}

		urls := make([]string, 0, len(comics))

		for _, comic := range comics {
			urls = append(urls, comic.ImgURL)
		}

		err := writeJson(w, urls)

		if err != nil {
			return errors.Join(ErrInternal, err)
		}

		return nil
	}
}
