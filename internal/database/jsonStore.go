package database

import (
	"os"
	"yadro-go-course/internal/comic"
	"yadro-go-course/pkg/jsonl"
)

type JsonDB struct {
	file *os.File
}

func NewJsonDB(file *os.File) *JsonDB {
	return &JsonDB{
		file: file,
	}
}

func (db *JsonDB) SaveAll(comics map[int]*comic.Comic) error {
	if err := db.file.Truncate(0); err != nil {
		return err
	}

	if _, err := db.file.Seek(0, 0); err != nil {
		return err
	}

	wr := jsonl.NewWriter(db.file)
	for _, c := range comics {
		if err := wr.WriteJson(c); err != nil {
			return err
		}
	}

	return wr.Flush()
}
func (db *JsonDB) Save(c *comic.Comic) error {
	wr := jsonl.NewWriter(db.file)
	if err := wr.WriteJson(c); err != nil {
		return err
	}

	return wr.Flush()
}
func (db *JsonDB) Read() (map[int]*comic.Comic, error) {
	comics := make(map[int]*comic.Comic)
	sc := jsonl.NewScanner(db.file)

	for sc.Scan() {
		var readComic comic.Comic

		if err := sc.Json(&readComic); err != nil {
			return nil, err
		}

		comics[readComic.ID] = &readComic
	}

	return comics, nil
}
