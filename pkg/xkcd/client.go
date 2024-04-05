package xkcd

import (
	"context"
	"yadro-go-course/pkg/words"
)

type Client struct {
	fetcher *fetcher
	stemmer *words.Stemmer
}

func NewClient(source string) *Client {
	return &Client{
		fetcher: newFetcher(source),
		stemmer: words.NewStemmer(nil),
	}
}

func (c *Client) GetFirst(ctx context.Context, num int) map[int]*Comic {
	fetched := c.fetcher.getFirst(ctx, num)
	comics := make(map[int]*Comic, len(fetched))

	for _, dto := range fetched {
		comic := dto.toComic(c.stemmer)
		if comic != nil {
			comics[comic.ID] = comic
		}
	}

	return comics
}
