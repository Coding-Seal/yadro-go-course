package main

import (
	"bufio"
	"github.com/kljensen/snowball"
	"io"
)

type Stemmer struct {
	stopWords map[string]struct{}
}

// NewStemmer is a function to construct new Stemmer.
// If stopWords == nil uses default dictionary https://www.ranks.nl/stopwords
func NewStemmer(stopWords map[string]struct{}) *Stemmer {
	if stopWords != nil {
		return &Stemmer{stopWords: stopWords}
	}
	return &Stemmer{stopWords: map[string]struct{}{
		"a": {}, "about": {}, "above": {}, "after": {}, "again": {}, "against": {}, "all": {},
		"am": {}, "an": {}, "and": {}, "any": {}, "are": {}, "aren't": {}, "as": {}, "at": {},
		"be": {}, "because": {}, "been": {}, "before": {}, "being": {}, "below": {}, "between": {},
		"both": {}, "but": {}, "by": {}, "can't": {}, "cannot": {}, "could": {}, "couldn't": {},
		"did": {}, "didn't": {}, "do": {}, "does": {}, "doesn't": {}, "doing": {}, "don't": {},
		"down": {}, "during": {}, "each": {}, "few": {}, "for": {}, "from": {}, "further": {},
		"had": {}, "hadn't": {}, "has": {}, "hasn't": {}, "have": {}, "haven't": {}, "having": {},
		"he": {}, "he'd": {}, "he'll": {}, "he's": {}, "her": {}, "here": {}, "here's": {}, "hers": {},
		"herself": {}, "him": {}, "himself": {}, "his": {}, "how": {}, "how's": {}, "i": {}, "i'd": {},
		"i'll": {}, "i'm": {}, "i've": {}, "if": {}, "in": {}, "into": {}, "is": {}, "isn't": {},
		"it": {}, "it's": {}, "its": {}, "itself": {}, "let's": {}, "me": {}, "more": {}, "most": {},
		"mustn't": {}, "my": {}, "myself": {}, "no": {}, "nor": {}, "not": {}, "of": {}, "off": {},
		"on": {}, "once": {}, "only": {}, "or": {}, "other": {}, "ought": {}, "our": {}, "ours": {},
		"ourselves": {}, "out": {}, "over": {}, "own": {}, "same": {}, "shan't": {}, "she": {},
		"she'd": {}, "she'll": {}, "she's": {}, "should": {}, "shouldn't": {}, "so": {}, "some": {},
		"such": {}, "than": {}, "that": {}, "that's": {}, "the": {}, "their": {}, "theirs": {},
		"them": {}, "themselves": {}, "then": {}, "there": {}, "there's": {}, "these": {}, "they": {},
		"they'd": {}, "they'll": {}, "they're": {}, "they've": {}, "this": {}, "those": {},
		"through": {}, "to": {}, "too": {}, "under": {}, "until": {}, "up": {}, "very": {}, "was": {},
		"wasn't": {}, "we": {}, "we'd": {}, "we'll": {}, "we're": {}, "we've": {}, "were": {},
		"weren't": {}, "what": {}, "what's": {}, "when": {}, "when's": {}, "where": {}, "where's": {},
		"which": {}, "while": {}, "who": {}, "who's": {}, "whom": {}, "why": {}, "why's": {},
		"with": {}, "won't": {}, "would": {}, "wouldn't": {}, "you": {}, "you'd": {}, "you'll": {},
		"you're": {}, "you've": {}, "your": {}, "yours": {}, "yourself": {}, "yourselves": {}}}

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
