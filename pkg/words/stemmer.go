package words

import (
	"bufio"
	"github.com/kljensen/snowball"
	"io"
	"regexp"
	"slices"
	"strings"
)

var nonWordSymbolRegexp = regexp.MustCompile("[^0-9A-Za-z_]+")

type Stemmer struct {
	stopWords  map[string]struct{}
	leaveWords int
}

// NewStemmer is a function to construct new Stemmer.
// If stopWords == nil uses default dictionary https://www.ranks.nl/stopwords
func NewStemmer(stopWords map[string]struct{}, leaveWords int) *Stemmer {
	if stopWords != nil {
		return &Stemmer{stopWords: stopWords, leaveWords: leaveWords}
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
		"you're": {}, "you've": {}, "your": {}, "yours": {}, "yourself": {}, "yourselves": {},
		"alt": {}, "text": {}, "title": {}}, leaveWords: leaveWords}
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

func ParsePhrase(phrase string) []string {
	phrase = nonWordSymbolRegexp.ReplaceAllString(phrase, " ")
	words := strings.Fields(phrase)
	deleted := 0

	for i := 0; i < len(words)-deleted; i++ {
		word := strings.ToLower(words[i])
		if word != "" {
			words[i] = word
		} else {
			words[i] = words[len(words)-deleted-1]
			i--
			deleted++
		}
	}

	words = slices.Delete(words, len(words)-deleted, len(words))

	return words
}

// isStopWord is method, that check if word is not significant
func (s *Stemmer) isStopWord(word string) bool {
	_, ok := s.stopWords[word]
	return ok
}
func (s *Stemmer) Stem(words []string) []string {
	// using map to avoid duplicates
	stemmed := make(map[string]int)

	for _, word := range words {
		if len(word) < 4 || s.isStopWord(word) {
			continue
		} else {
			word, _ = snowball.Stem(word, "english", false)
			if len(word) < 3 {
				continue
			}

			stemmed[word] += 1
		}
	}

	type pair struct {
		word  string
		count int
	}
	// transform map into slice
	pairs := make([]pair, 0, len(stemmed))

	for k, v := range stemmed {
		pairs = append(pairs, pair{word: k, count: v})
	}

	slices.SortFunc(pairs, func(a, b pair) int {
		if b.count == a.count {
			if b.word < a.word {
				return 1
			} else if b.word > a.word {
				return -1
			} else {
				return 0
			}
		}

		return b.count - a.count
	})

	res := make([]string, min(s.leaveWords, len(pairs)))

	for i := 0; i < min(s.leaveWords, len(pairs)); i++ {
		res[i] = pairs[i].word
	}

	return res
}
