package xkcd

import (
	"context"
	"fmt"
	"testing"
	"time"
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
				b.Errorf(fmt.Sprintf("expected %d +- %d comics, got %d", lastID-1, EPS, len(comics)))
			}
		})
	}
}

func TestFetcher_SearchLastID(t *testing.T) {
	fetcher := NewFetcher("https://xkcd.com", 0)
	lastID, err := fetcher.LastID(context.Background())

	if err != nil {
		t.Errorf("failed to get last id (lifehack) : %v", err)
	}

	id, err := fetcher.SearchLastID(context.Background())
	if err != nil {
		t.Errorf("failed to search last id : %v", err)
	}

	if id != lastID {
		t.Errorf("got %d, expected %d", id, lastID)
	}
}
