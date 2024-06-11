package handlers

import (
	"errors"
	"net/http"

	"yadro-go-course/pkg/http-util"

	"yadro-go-course/internal/core/services"
)

func Update(fetcherSrv *services.Fetcher) http_util.ErrHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := fetcherSrv.Update(r.Context())
		if err != nil {
			return errors.Join(http_util.ErrInternal, err)
		}

		return nil
	}
}

func Search(searchSrv *services.Search) http_util.ErrHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		phrase := r.URL.Query().Get("search")
		comics := searchSrv.SearchComics(r.Context(), phrase, 10)

		if len(comics) == 0 {
			return http_util.ErrNotFound
		}

		urls := make([]string, 0, len(comics))

		for _, comic := range comics {
			urls = append(urls, comic.ImgURL)
		}

		err := http_util.WriteJson(w, urls)
		if err != nil {
			return errors.Join(http_util.ErrInternal, err)
		}

		return nil
	}
}
