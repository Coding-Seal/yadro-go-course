package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var input string
	flag.StringVar(&input, "s", "", "String to stem")
	var stopWordsFileName string
	flag.StringVar(&stopWordsFileName, "file", "", "File with stop words")
	flag.Parse()
	if input == "" {
		fmt.Println("Provide string to stem using -s flag")
		os.Exit(1)
	}
	// Handle file with stopWords
	var stopWords map[string]struct{}
	if stopWordsFileName != "" {
		stopWordsFile, err := os.Open(stopWordsFileName)
		if err != nil {
			fmt.Printf("Could not open file \"%s\"\n", stopWordsFileName)
			os.Exit(1)
		}
		stopWords = ParseStopWords(stopWordsFile)
	}

	words := strings.Fields(input)
	stemmer := NewStemmer(stopWords)
	stemmed := stemmer.Stem(words)
	fmt.Println(strings.Join(stemmed, " "))
}
