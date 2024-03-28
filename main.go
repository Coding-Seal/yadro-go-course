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
	var stopWordsFile string
	flag.StringVar(&stopWordsFile, "file", "", "File with stop words")
	flag.Parse()
	if input == "" {
		fmt.Println("Provide string to stem using -s flag")
		os.Exit(1)
	}
	if stopWordsFile == "" {
		fmt.Println("Provide file with stop words using --file flag")
		os.Exit(1)
	}
	words := strings.Fields(strings.ToLower(input))
	stopWords, err := os.Open(stopWordsFile)
	if err != nil {
		fmt.Printf("Could not open file \"%s\"\n", stopWordsFile)
		os.Exit(1)
	}
	stemmer := NewStemmer(ParseStopWords(stopWords))
	stemmed := stemmer.Stem(words)
	fmt.Println(strings.Join(stemmed, " "))
}
