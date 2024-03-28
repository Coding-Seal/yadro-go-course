package main

import (
	"bufio"
	"github.com/kljensen/snowball"
	"io"
)

type Stemmer struct {
	stopWords map[string]struct{}
}

func NewStemmer(stopWords map[string]struct{}) *Stemmer {
	return &Stemmer{stopWords: stopWords}
}

func ParseStopWords(reader io.Reader) map[string]struct{} {
	stopWords := make(map[string]struct{})
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		if word != "" {
			stopWords[word] = struct{}{}
		}
	}
	return stopWords
}

// isStopWord is method, that check if word is not significant
func (s *Stemmer) isStopWord(word string) bool {
	_, ok := s.stopWords[word]
	return ok
}
func (s *Stemmer) Stem(words []string) []string {
	// using map to avoid duplicates
	stemmed := make(map[string]struct{})
	for _, word := range words {
		if s.isStopWord(word) || word == "" {
			continue
		} else {
			word, _ = snowball.Stem(word, "english", false)
			stemmed[word] = struct{}{}
		}
	}
	// transform map into slice
	keys := make([]string, 0, len(stemmed))
	for k := range stemmed {
		keys = append(keys, k)
	}
	return keys
}
