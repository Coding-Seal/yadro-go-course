package main

import (
	"github.com/kljensen/snowball"
)

// isStopWord is function, that check if word is not significant
// list of words from https://www.ranks.nl/stopwords
func isStopWord(word string) bool {
	switch word {
	case "a", "about", "above", "after", "again", "against", "all", "am", "an", "and", "any",
		"are", "aren't", "as", "at", "be", "because", "been", "before", "being", "below", "between",
		"both", "but", "by", "can't", "cannot", "could", "couldn't", "did", "didn't", "do", "does",
		"doesn't", "doing", "don't", "down", "during", "each", "few", "for", "from", "further", "had",
		"hadn't", "has", "hasn't", "have", "haven't", "having", "he", "he'd", "he'll", "he's", "her",
		"here", "here's", "hers", "herself", "him", "himself", "his", "how", "how's", "i", "i'd",
		"i'll", "i'm", "i've", "if", "in", "into", "is", "isn't", "it", "it's", "its", "itself",
		"let's", "me", "more", "most", "mustn't", "my", "myself", "no", "nor", "not", "of", "off",
		"on", "once", "only", "or", "other", "ought", "our", "ours", "ourselves", "out", "over",
		"own", "same", "shan't", "she", "she'd", "she'll", "she's", "should", "shouldn't", "so",
		"some", "such", "than", "that", "that's", "the", "their", "theirs", "them", "themselves",
		"then", "there", "there's", "these", "they", "they'd", "they'll", "they're", "they've", "this",
		"those", "through", "to", "too", "under", "until", "up", "very", "was", "wasn't", "we", "we'd",
		"we'll", "we're", "we've", "were", "weren't", "what", "what's", "when", "when's", "where",
		"where's", "which", "while", "who", "who's", "whom", "why", "why's", "with", "won't", "would",
		"wouldn't", "you", "you'd", "you'll", "you're", "you've", "your", "yours", "yourself", "yourselves":
		return true
	default:
		return false

	}
}
func Stem(words []string) []string {
	stemmed := make([]string, 0, len(words)/2)
	for _, word := range words {
		if isStopWord(word) || word == "" {
			continue
		} else {
			word, _ = snowball.Stem(word, "english", false)
			stemmed = append(stemmed, word)
		}
	}
	return stemmed
}
