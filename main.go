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
	flag.Parse()
	if input == "" {
		fmt.Println("Provide string to stem using -s flag")
		os.Exit(1)
	}
	words := strings.Fields(strings.ToLower(input))
	stemmed := Stem(words)
	fmt.Println(strings.Join(stemmed, " "))
}
