package words

import (
	"bytes"
	"maps"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStemmer(t *testing.T) {
	if NewStemmer(nil) == nil {
		t.Error("NewStemmer should constructed")
	}

	if NewStemmer(map[string]struct{}{"word": {}, "someWord": {}}) == nil {
		t.Error("NewStemmer should constructed")
	}
}

func TestParseStopWords(t *testing.T) {
	buf := bytes.NewBufferString("an apple a day")
	stopWords := map[string]struct{}{"an": {}, "apple": {}, "a": {}, "day": {}}
	words := ParseStopWords(buf)

	if !maps.Equal(stopWords, words) {
		t.Error("some words don't match")
	}
}

func TestParsePhrase(t *testing.T) {
	type words struct {
		phrase string
		words  []string
	}

	tests := []words{
		{"So close, no matter how far", []string{"so", "close", "no", "matter", "how", "far"}},
		{"Couldn't^be	much&more%from the heart", []string{"couldn", "t", "be", "much", "more", "from", "the", "heart"}},
		{"12 34 56 78", []string{"12", "34", "56", "78"}},
	}
	for _, test := range tests {
		res := ParsePhrase(test.phrase)
		expected := test.words
		sort.Strings(expected)
		sort.Strings(res)
		assert.Equal(t, expected, res)
	}
}

func TestStemmer(t *testing.T) {
	stemmer := NewStemmer(nil)
	stemmer.Stem(ParsePhrase("Forever trusting who we are"))
}
