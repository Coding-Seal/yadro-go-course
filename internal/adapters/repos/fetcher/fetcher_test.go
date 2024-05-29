package fetcher

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	testfetcher "yadro-go-course/test/fetcher"
)

func TestFetcher_LastComicID(t *testing.T) {
	lastID := 20
	srv := testfetcher.NewMockXKCD(lastID)
	t.Cleanup(srv.Close)
	fetcher := NewFetcher(srv.URL, 10)
	_, err := fetcher.LastComicID(context.Background())
	assert.NoError(t, err)
}

func TestFetcher_Comics(t *testing.T) {
	ctx := context.Background()
	lastID := 20
	srv := testfetcher.NewMockXKCD(lastID)
	t.Cleanup(srv.Close)
	fetcher := NewFetcher(srv.URL, 10)
	id, err := fetcher.LastComicID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, lastID, id)

	ids, comics := fetcher.Comics(ctx, id)
	for i := 1; i <= id; i++ {
		ids <- i
	}

	var result int
	for range comics {
		result++
	}
}
