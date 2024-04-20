package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"yadro-go-course/internal/comic"
)

type JsonDB struct {
	file *os.File
}

func NewJsonDB(file *os.File) *JsonDB {
	return &JsonDB{
		file: file,
	}
}

func (db *JsonDB) Save(comics map[int]*comic.Comic) {
	_ = db.file.Truncate(0)
	_, _ = db.file.Seek(0, 0)
	encoder := json.NewEncoder(db.file)
	err := encoder.Encode(comics)

	if err != nil {
		log.Println(fmt.Errorf("error while saving comics: %w", err))
	}
}
func (db *JsonDB) Read() map[int]*comic.Comic {
	comics := make(map[int]*comic.Comic)
	decoder := json.NewDecoder(db.file)
	err := decoder.Decode(&comics)

	if err != nil {
		if errors.Is(err, io.EOF) {
			return comics
		}

		log.Println(fmt.Errorf("error while loading comics: %w", err))
	}

	return comics
}
