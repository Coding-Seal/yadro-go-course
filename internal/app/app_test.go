package app

import (
	"os"
	"testing"
	"yadro-go-course/pkg/words"
)

var phrase = "I'm following your questions"
var sourceURL = "https://xkcd.com"
var dbPath = "../../somefile.jsonl"
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
	if err := client.LoadComics(); err != nil {
		b.Error("Could not parse comics", err)
		b.FailNow()
	}

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
	if err := client.LoadComics(); err != nil {
		b.Error("Could not parse comics", err)
		b.FailNow()
	}

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
	if err := client.LoadComics(); err != nil {
		b.Error("Could not parse comics", err)
		b.FailNow()
	}

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

	if err := client.LoadComics(); err != nil {
		b.Error("Could not parse comics", err)
	}
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
	if err := client.LoadComics(); err != nil {
		t.Error("Could not parse comics", err)
		t.FailNow()
	}

	client.BuildIndex()

	indexFound := client.SearchIndex(phrase)
	mapFound := client.SearchComics(phrase)

	for foundComic, score := range mapFound {
		sc, ok := indexFound[foundComic]

		if !ok {
			t.Errorf("No comic %v in %v", foundComic, indexFound)
		}

		if sc != score {
			t.Errorf("score should be the same %d != %d, comic %v", sc, score, foundComic)
		}
	}

	for foundComic, score := range indexFound {
		sc, ok := mapFound[foundComic]

		if !ok {
			t.Errorf("No comic %v in %v", foundComic, mapFound)
		}

		if sc != score {
			t.Errorf("score should be the same %d != %d, comic %v", sc, score, foundComic)
		}
	}
}
