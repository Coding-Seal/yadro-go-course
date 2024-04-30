package fetcher

import (
	"context"
	"errors"
	"yadro-go-course/internal/core/models"
	"yadro-go-course/internal/core/ports"
	"yadro-go-course/pkg/xkcd"
)

var ErrInternal = errors.New("internal error")

type Fetcher struct {
	fetcher xkcd.Fetcher
}

var _ ports.ComicFetcherRepo = (*Fetcher)(nil)

func (f *Fetcher) LastComicID(ctx context.Context) (int, error) {
	id, err := f.fetcher.LastID(ctx)
	if err != nil {
		return 0, errors.Join(ErrInternal, err)
	}
	return id, nil
}

func (f *Fetcher) Comics(ctx context.Context, limit int) (chan<- int, <-chan models.Comic) {
	jobsCh, fetchedCh := f.fetcher.Comics(ctx, limit)
	comicsCh := make(chan models.Comic, limit)
	go func() {
		for i := 0; i < limit; i++ {
			fetchedComic, ok := <-fetchedCh
			if !ok {
				break
			}
			if fetchedComic.Err() != nil {
				continue
			}
			comicsCh <- models.Comic{
				ID:               fetchedComic.Comic.ID,
				Title:            fetchedComic.Comic.Title,
				Date:             fetchedComic.Comic.Date,
				ImgURL:           fetchedComic.Comic.ImgURL,
				News:             fetchedComic.Comic.News,
				SafeTitle:        fetchedComic.Comic.SafeTitle,
				Transcription:    fetchedComic.Comic.Transcription,
				AltTranscription: fetchedComic.Comic.AltTranscription,
				Link:             fetchedComic.Comic.Link,
			}
		}
	}()
	return jobsCh, comicsCh
}
