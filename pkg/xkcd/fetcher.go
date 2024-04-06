package xkcd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
	"yadro-go-course/pkg/comic"
	"yadro-go-course/pkg/words"
)

type Fetcher struct {
	client *http.Client
	source string
}

func NewFetcher(source string) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   time.Second * 7,
		},
		source: source,
	}
}
func (f *Fetcher) GetComics(ctx context.Context, ids []int) map[int]*FetchedComic {
	wg := &sync.WaitGroup{}
	mu := sync.Mutex{}
	comics := make(map[int]*FetchedComic, len(ids))

	for id := range ids {
		wg.Add(1)

		go func(id int) {
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
func (f *Fetcher) Get(ctx context.Context, id int) *FetchedComic {
	req, err := http.NewRequestWithContext(ctx, "GET", buildURL(f.source, id), nil)
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

func buildURL(source string, id int) string {
	return fmt.Sprintf("%s/%d/info.0.json", source, id)
}

type FetchedComic struct {
	ID               int    `json:"num"`
	ImgURL           string `json:"img"`
	Title            string `json:"title"`
	Transcription    string `json:"transcript"`
	AltTranscription string `json:"alt"`
}

func (d *FetchedComic) ToComic(stemmer *words.Stemmer) *comic.Comic {
	return &comic.Comic{
		ID:       d.ID,
		ImgURL:   d.ImgURL,
		Keywords: stemmer.Stem(words.ParsePhrase(d.AltTranscription + " " + d.Transcription + " " + d.Title)),
	}
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
