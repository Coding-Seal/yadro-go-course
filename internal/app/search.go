package app

import (
	"yadro-go-course/internal/core/models"
	"yadro-go-course/pkg/words"
	"yadro-go-course/pkg/xkcd"
)

func giveWordScore(inPhrase, inKeywords int) int {
	return min(inPhrase*inKeywords, 3)
}

func (a *App) toComic(c *xkcd.Comic) *models.Comic {
	return &models.Comic{
		ID:       c.ID,
		Title:    c.Title,
		Date:     c.Date,
		ImgURL:   c.ImgURL,
		Keywords: a.stemmer.Stem(words.ParsePhrase(c.Title + " " + c.AltTranscription + " " + c.Transcription + " " + c.News)),
	}
}
func (a *App) SearchComics(searchPhrase string) map[*models.Comic]int {
	searchWords := a.stemmer.Stem(words.ParsePhrase(searchPhrase))
	res := make(map[*models.Comic]int, 0)

	for _, c := range a.comics {
		score := 0
		for word, numWords := range searchWords {
			score += giveWordScore(numWords, c.Keywords[word])
		}

		if score > 0 {
			res[c] = score
		}
	}

	return res
}
func (a *App) BuildIndex() {
	a.index = make(map[string][]int, len(a.comics)*10)
	for id, c := range a.comics {
		for word := range c.Keywords {
			a.index[word] = append(a.index[word], id)
		}
	}
}
func (a *App) SearchIndex(searchPhrase string) map[*models.Comic]int {
	searchWords := a.stemmer.Stem(words.ParsePhrase(searchPhrase))
	foundComics := make(map[*models.Comic]int)

	for word, num := range searchWords {
		for _, id := range a.index[word] {
			foundComics[a.comics[id]] += giveWordScore(num, a.comics[id].Keywords[word])
		}
	}

	return foundComics
}
