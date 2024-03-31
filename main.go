package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var input string

	var stopWordsFileName string

	flag.StringVar(&input, "s", "", "String to stem")
	flag.StringVar(&stopWordsFileName, "file", "", "File with stop words")
	flag.Parse()

	if input == "" {
		fmt.Println("Provide string to stem using -s flag")
		os.Exit(1)
	}

	words := ParsePhrase(input)

	var stopWords map[string]struct{}
	// Handle optional file with stopWords
	if stopWordsFileName != "" {
		stopWordsFile, err := os.Open(stopWordsFileName)
		if err != nil {
			fmt.Printf("Could not open file \"%s\"\n", stopWordsFileName)
			os.Exit(1)
		}

		stopWords = ParseStopWords(stopWordsFile)
	}

	stemmer := NewStemmer(stopWords)
	stemmed := stemmer.Stem(words)
	fmt.Println(strings.Join(stemmed, " "))
}
