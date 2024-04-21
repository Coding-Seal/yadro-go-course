package app

import (
	"os"
	"slices"
	"testing"
	"yadro-go-course/pkg/words"
)

var phrase = "I'm following your questions"
var sourceURL = "https://xkcd.com"
var dbPath = "../../somefile.json"
var stopWordsPath = "../../stopwords.txt"
var parallelLimit = 100

func BenchmarkApp_SearchComics(b *testing.B) {
	dbFile, err := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		b.Error("Could not open db file", err)
		b.FailNow()
	}

	defer dbFile.Close()

	stopWordsFile, err := os.Open(stopWordsPath)

	if err != nil {
		b.Error("Could not open stop words file", err)
		b.FailNow()
	}

	defer stopWordsFile.Close()
	stopWordsMap := words.ParseStopWords(stopWordsFile)

	client := NewApp(sourceURL, dbFile, stopWordsMap, parallelLimit)
	client.LoadComics()
	b.ResetTimer()
	client.SearchComics(phrase)
}

func BenchmarkApp_SearchIndex(b *testing.B) {
	dbFile, err := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		b.Error("Could not open db file", err)
		b.FailNow()
	}

	defer dbFile.Close()

	stopWordsFile, err := os.Open(stopWordsPath)

	if err != nil {
		b.Error("Could not open stop words file", err)
		b.FailNow()
	}

	defer stopWordsFile.Close()
	stopWordsMap := words.ParseStopWords(stopWordsFile)

	client := NewApp(sourceURL, dbFile, stopWordsMap, parallelLimit)
	client.LoadComics()

	client.BuildIndex()
	b.ResetTimer()
	client.SearchIndex(phrase)
}
func BenchmarkApp_BuildIndex(b *testing.B) {
	dbFile, err := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		b.Error("Could not open db file", err)
		b.FailNow()
	}

	defer dbFile.Close()

	stopWordsFile, err := os.Open(stopWordsPath)

	if err != nil {
		b.Error("Could not open stop words file", err)
		b.FailNow()
	}

	defer stopWordsFile.Close()
	stopWordsMap := words.ParseStopWords(stopWordsFile)

	client := NewApp(sourceURL, dbFile, stopWordsMap, parallelLimit)
	client.LoadComics()
	b.ResetTimer()
	client.BuildIndex()
}
func BenchmarkApp_LoadComics(b *testing.B) {
	dbFile, err := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		b.Error("Could not open db file", err)
		b.FailNow()
	}

	b.StopTimer()

	defer dbFile.Close()

	stopWordsFile, err := os.Open(stopWordsPath)

	if err != nil {
		b.Error("Could not open stop words file", err)
		b.FailNow()
	}

	defer stopWordsFile.Close()
	stopWordsMap := words.ParseStopWords(stopWordsFile)

	client := NewApp(sourceURL, dbFile, stopWordsMap, parallelLimit)

	b.StartTimer()
	client.LoadComics()
}

func TestApp_SearchComics(t *testing.T) {
	dbFile, err := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		t.Error("Could not open db file", err)
		t.FailNow()
	}

	defer dbFile.Close()

	stopWordsFile, err := os.Open(stopWordsPath)

	if err != nil {
		t.Error("Could not open stop words file", err)
		t.FailNow()
	}

	defer stopWordsFile.Close()
	stopWordsMap := words.ParseStopWords(stopWordsFile)

	client := NewApp(sourceURL, dbFile, stopWordsMap, parallelLimit)
	client.LoadComics()

	client.BuildIndex()

	indexFound := client.SearchIndex(phrase)
	mapFound := client.SearchComics(phrase)

	for _, foundComic := range mapFound {
		if !slices.Contains(indexFound, foundComic) {
			t.Errorf("Not comics %v in %v", foundComic, indexFound)
		}
	}

	for _, foundComic := range indexFound {
		if !slices.Contains(mapFound, foundComic) {
			t.Errorf("Not comics %v in %v", foundComic, mapFound)
		}
	}
}
