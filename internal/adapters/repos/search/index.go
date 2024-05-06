package search

import (
	"context"
	"sync"
	"yadro-go-course/internal/core/models"
	"yadro-go-course/internal/core/ports"
	"yadro-go-course/pkg/words"
)

type Index struct {
	ind     map[string]map[int]struct{}
	mu      sync.RWMutex
	stemmer *words.Stemmer
}

var _ ports.SearchComicsRepo = (*Index)(nil)

func NewIndex(stemmer *words.Stemmer) *Index {
	return &Index{
		ind:     make(map[string]map[int]struct{}),
		stemmer: stemmer,
	}
}

func (index *Index) SearchComics(ctx context.Context, query string) map[int]int {
	parsed := words.ParsePhrase(query)
	stemmed := index.stemmer.Stem(parsed)
	found := make(map[int]int)

	for word := range stemmed {
		index.mu.RLock()
		for id := range index.ind[word] {
			found[id]++
		}
		index.mu.RUnlock()
	}

	return found
}
func (index *Index) AddComic(ctx context.Context, comic models.Comic) {
	keywords := index.stemmer.Stem(words.ParsePhrase(comic.Title + " " + comic.SafeTitle +
		" " + comic.Transcription + " " + comic.AltTranscription))
	for word := range keywords {
		index.mu.Lock()
		if _, ok := index.ind[word]; !ok {
			index.ind[word] = make(map[int]struct{})
		}

		index.ind[word][comic.ID] = struct{}{}
		index.mu.Unlock()
	}
}
func (index *Index) Build(ctx context.Context, repo ports.ComicsRepo) error {
	comics, err := repo.ComicsAll(ctx)
	if err != nil {
		return err
	}

	for _, comic := range comics {
		index.AddComic(ctx, comic)
	}

	return nil
}
