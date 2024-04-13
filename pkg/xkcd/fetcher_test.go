package xkcd

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var concurrencyLimits = []int{5, 10, 25, 50, 100, 200, 225, 250, 275, 258, 300, 350, 400, 450, 500, 750, 1000, 1500}

const EPS = 0

func BenchmarkFetcher_GetAllComics(b *testing.B) {
	fetcher := NewFetcher("https://xkcd.com", 0)
	lastID, err := fetcher.GetLastID(context.Background())
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for _, limit := range concurrencyLimits {
		fetcher.concurrencyLimit = limit

		b.Run(fmt.Sprintf("cuncurrency_limit_%d", limit), func(b *testing.B) {
			b.N = 100
			b.ResetTimer()
			comics := fetcher.GetAllComics(context.Background(), lastID)
			b.StopTimer()
			time.Sleep(1 * time.Second)
			if len(comics)+EPS < lastID-1 {
				b.Errorf(fmt.Sprintf("expected %d +- %d comics, got %d", lastID-1, EPS, len(comics)))
			}
		})
	}
}

func BenchmarkFetcher_GetComics(b *testing.B) {
	fetcher := NewFetcher("https://xkcd.com", 0)
	comicsToFetch := make([]int, 0, 100)
	for i := 1; i <= 100; i++ {
		comicsToFetch = append(comicsToFetch, i)
	}

	b.ResetTimer()
	for _, limit := range concurrencyLimits {
		fetcher.concurrencyLimit = limit

		b.Run(fmt.Sprintf("cuncurrency_limit_%d", limit), func(b *testing.B) {
			b.N = 100
			b.ResetTimer()
			comics := fetcher.GetComics(context.Background(), comicsToFetch)
			b.StopTimer()
			time.Sleep(1 * time.Second)
			if len(comics)+EPS < len(comicsToFetch) {
				b.Errorf(fmt.Sprintf("expected %d +- %d comics, got %d", len(comicsToFetch), EPS, len(comics)))
			}
		})
	}
}
