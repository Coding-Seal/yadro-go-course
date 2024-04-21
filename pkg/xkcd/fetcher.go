package xkcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

var (
	errNotFound = errors.New(http.StatusText(http.StatusNotFound))
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

func (f *Fetcher) Comics(ctx context.Context, numComics int) (chan<- int, <-chan FetchedComic) {
	jobs := make(chan int, numComics)
	results := make(chan FetchedComic, numComics)

	for w := 0; w <= f.concurrencyLimit; w++ {
		go func() {
			for id := range jobs {
				comic, err := f.Comic(ctx, id)
				results <- FetchedComic{comic, err}
			}
		}()
	}

	return jobs, results
}

func (f *Fetcher) Comic(ctx context.Context, id int) (*Comic, error) {
	url := fmt.Sprintf("%s/%d/info.0.json", f.source, id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", `application/json`)
	resp, err := f.client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(http.StatusText(resp.StatusCode))
	}
	defer resp.Body.Close()

	return parseJsonComic(resp.Body), nil
}

func (f *Fetcher) LastID(ctx context.Context) (int, error) {
	url := fmt.Sprintf("%s/info.0.json", f.source)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return 0, fmt.Errorf("LastID : %w", err)
	}

	req.Header.Add("Accept", `application/json`)
	resp, err := f.client.Do(req)

	if err != nil {
		return 0, fmt.Errorf("LastID : %w", err)
	}

	defer resp.Body.Close()
	fetched := parseJsonComic(resp.Body)

	if fetched == nil {
		return 0, errors.New("could not parse last comic")
	}

	return fetched.ID, nil
}

func (f *Fetcher) isComicPresent(ctx context.Context, id int) (bool, error) {
	_, err := f.Comic(ctx, id)

	if err != nil {
		if errors.Is(err, errNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}
func (f *Fetcher) SearchLastID(ctx context.Context) (int, error) {
	left, right := 1, math.MaxInt
	leftPresent, err := f.isComicPresent(ctx, left)

	if err != nil {
		return 0, err
	}

	rightPresent, err := f.isComicPresent(ctx, right)

	if err != nil {
		return 0, err
	}

	for left+1 < right && leftPresent && !rightPresent {
		pivot := left + (right-left)/2
		pivotPresent, err := f.isComicPresent(ctx, pivot)

		if err != nil {
			return 0, err
		}

		if pivotPresent {
			left = pivot
		} else {
			right = pivot
		}
	}

	return left, nil
}

type ParsedComic struct {
	Day              int        `json:"day,string"`
	Month            time.Month `json:"month,string"`
	Year             int        `json:"year,string"`
	ID               int        `json:"num"`
	News             string     `json:"news"`
	SafeTitle        string     `json:"safe_title"`
	ImgURL           string     `json:"img"`
	Title            string     `json:"title"`
	Transcription    string     `json:"transcript"`
	AltTranscription string     `json:"alt"`
	Link             string     `json:"link"`
}

func (c *ParsedComic) toComic() *Comic {
	return &Comic{
		ID:               c.ID,
		Date:             time.Date(c.Year, c.Month, c.Day, 0, 0, 0, 0, time.UTC),
		News:             c.News,
		SafeTitle:        c.SafeTitle,
		ImgURL:           c.ImgURL,
		Title:            c.Title,
		Transcription:    c.Transcription,
		AltTranscription: c.AltTranscription,
		Link:             c.Link,
	}
}

func parseJsonComic(r io.Reader) *Comic {
	var dto ParsedComic

	decoder := json.NewDecoder(r)
	err := decoder.Decode(&dto)

	if err != nil {
		return nil
	}

	return dto.toComic()
}

type FetchedComic struct {
	Comic *Comic
	err   error
}

func (c FetchedComic) Err() error {
	return c.err
}
