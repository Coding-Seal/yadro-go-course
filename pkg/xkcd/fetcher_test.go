package xkcd

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	testfetcher "yadro-go-course/test/fetcher"
)

var concurrencyLimits = []int{50, 100, 200, 250, 300, 350, 400, 500, 750, 1000, 1500}

const EPS = 0

func BenchmarkFetcher_Comics(b *testing.B) {
	fetcher := NewFetcher("https://xkcd.com", 0)

	lastID, err := fetcher.LastID(context.Background())
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	numComics := 0

	for _, limit := range concurrencyLimits {
		fetcher.concurrencyLimit = limit

		b.Run(fmt.Sprintf("cuncurrency_limit_%d", limit), func(b *testing.B) {
			b.ResetTimer()

			ids, comics := fetcher.Comics(context.Background(), lastID)
			for id := 1; id <= lastID; id++ {
				ids <- id
			}
			close(ids)

			for i := 0; i < lastID; i++ {
				fetchedComic := <-comics
				if fetchedComic.Err() == nil {
					numComics++
				}
			}

			b.StopTimer()
			time.Sleep(1 * time.Second)

			if numComics+EPS < lastID-1 {
				b.Errorf("expected %d +- %d comics, got %d", lastID-1, EPS, len(comics))
			}
		})
	}
}

func TestFetcher_SearchLastID(t *testing.T) {
	lastID := 5
	srv := testfetcher.NewMockXKCD(lastID)
	t.Cleanup(srv.Close)
	fetcher := NewFetcher(srv.URL, 0)
	lastID, err := fetcher.LastID(context.Background())
	assert.NoError(t, err)
	id, err := fetcher.SearchLastID(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, lastID, id)
}

func TestFetcher_Comic(t *testing.T) {
	srv := testfetcher.NewMockXKCD(5)
	t.Cleanup(srv.Close)
	fetcher := NewFetcher(srv.URL, 0)
	c, err := fetcher.Comic(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, c.ID, 1)
}

func TestFetcher_Comics(t *testing.T) {
	lastID := 11
	srv := testfetcher.NewMockXKCD(lastID)
	t.Cleanup(srv.Close)

	fetcher := NewFetcher(srv.URL, 10)

	ids, fet := fetcher.Comics(context.Background(), 10)
	for i := 1; i <= 11; i++ {
		ids <- i
	}
	close(ids)

	for i := 1; i <= 11; i++ {
		f := <-fet
		assert.NoError(t, f.Err())
	}
}
