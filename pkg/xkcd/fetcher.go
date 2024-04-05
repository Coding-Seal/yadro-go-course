package xkcd

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"yadro-go-course/pkg/words"
)

type fetcher struct {
	client *http.Client
	source string
}

func newFetcher(source string) *fetcher {
	return &fetcher{
		client: http.DefaultClient,
		source: source,
	}
}
func (f *fetcher) getFirst(ctx context.Context, num int) map[int]*fetchedComic {
	wg := &sync.WaitGroup{}
	mu := sync.Mutex{}
	comics := make(map[int]*fetchedComic, num)

	for i := 0; i <= num; i++ {
		wg.Add(1)

		go func(id int) {
			comic := f.get(ctx, id)
			if comic != nil {
				mu.Lock()
				comics[id] = comic
				mu.Unlock()
			}

			wg.Done()
		}(i)
	}
	wg.Wait()

	return comics
}
func (f *fetcher) get(ctx context.Context, id int) *fetchedComic {
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
	urlBuilder := strings.Builder{}
	urlBuilder.WriteString(source)
	urlBuilder.WriteString("/")
	urlBuilder.WriteString(strconv.Itoa(id))
	urlBuilder.WriteString("/")
	urlBuilder.WriteString("info.0.json")

	return urlBuilder.String()
}

type fetchedComic struct {
	ID               int    `json:"num"`
	ImgURL           string `json:"img"`
	Title            string `json:"title"`
	Transcription    string `json:"transcript"`
	AltTranscription string `json:"alt"`
}

func (d *fetchedComic) toComic(stemmer *words.Stemmer) *Comic {
	return &Comic{
		ID:       d.ID,
		Title:    d.Title,
		ImgURL:   d.ImgURL,
		Keywords: stemmer.Stem(words.ParsePhrase(d.AltTranscription + " " + d.Transcription + " " + d.Title)),
	}
}
func parseJsonComic(r io.Reader) *fetchedComic {
	var dto fetchedComic

	decoder := json.NewDecoder(r)
	err := decoder.Decode(&dto)

	if err != nil {
		return nil
	}

	return &dto
}
