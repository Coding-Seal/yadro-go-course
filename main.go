package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	var input string
	flag.StringVar(&input, "s", "", "String to stem")
	flag.Parse()
	words := strings.Fields(strings.ToLower(input))
	stemmed := Stem(words)
	fmt.Println(strings.Join(stemmed, " "))
}
