package xkcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/sync/semaphore"
	"io"
	"net/http"
	"sync"
	"time"
)

type Fetcher struct {
	client           *http.Client
	source           string
	concurrencyLimit int
}

func NewFetcher(source string, concurrencyLimit int) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   time.Second * 7,
		},
		source:           source,
		concurrencyLimit: concurrencyLimit,
	}
}
func (f *Fetcher) GetComics(ctx context.Context, ids []int) map[int]*FetchedComic {
	wg := &sync.WaitGroup{}
	mu := sync.Mutex{}
	comics := make(map[int]*FetchedComic, len(ids))
	sem := semaphore.NewWeighted(int64(f.concurrencyLimit))

	for _, id := range ids {
		err := sem.Acquire(ctx, 1)
		if err != nil {
			break
		}

		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			defer sem.Release(1)

			comic := f.Get(ctx, id)

			if comic != nil {
				mu.Lock()
				comics[id] = comic
				mu.Unlock()
			}
		}(id)
	}

	wg.Wait()

	return comics
}

func (f *Fetcher) Get(ctx context.Context, id int) *FetchedComic {
	url := fmt.Sprintf("%s/%d/info.0.json", f.source, id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil
	}

	req.Header.Add("Accept", `application/json`)
	resp, err := f.client.Do(req)

	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	return parseJsonComic(resp.Body)
}

func (f *Fetcher) GetLastID(ctx context.Context) (int, error) {
	url := fmt.Sprintf("%s/info.0.json", f.source)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return 0, fmt.Errorf("GetLastID : %w", err)
	}

	req.Header.Add("Accept", `application/json`)
	resp, err := f.client.Do(req)

	if err != nil {
		return 0, fmt.Errorf("GetLastID : %w", err)
	}

	defer resp.Body.Close()
	fetched := parseJsonComic(resp.Body)

	if fetched == nil {
		return 0, errors.New("could not parse last comic")
	}

	return fetched.ID, nil
}

func (f *Fetcher) GetAllComics(ctx context.Context, lastID int) map[int]*FetchedComic {
	wg := &sync.WaitGroup{}
	mu := sync.Mutex{}
	sem := semaphore.NewWeighted(int64(f.concurrencyLimit))
	comics := make(map[int]*FetchedComic, lastID)

	for id := 1; id <= lastID; id++ {
		wg.Add(1)

		go func(id int) {
			err := sem.Acquire(ctx, 1)
			if err != nil {
				return
			}

			defer sem.Release(1)
			defer wg.Done()

			comic := f.Get(ctx, id)

			if comic != nil {
				mu.Lock()
				comics[id] = comic
				mu.Unlock()
			}
		}(id)
	}

	wg.Wait()

	return comics
}

type FetchedComic struct {
	ID               int    `json:"num"`
	ImgURL           string `json:"img"`
	Title            string `json:"title"`
	Transcription    string `json:"transcript"`
	AltTranscription string `json:"alt"`
}

func parseJsonComic(r io.Reader) *FetchedComic {
	var dto FetchedComic

	decoder := json.NewDecoder(r)
	err := decoder.Decode(&dto)

	if err != nil {
		return nil
	}

	return &dto
}
