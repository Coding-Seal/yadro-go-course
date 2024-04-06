package database

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"yadro-go-course/pkg/comic"
)

type JsonDB struct {
	rw io.ReadWriter
}

func NewJsonDB(rw io.ReadWriter) *JsonDB {
	return &JsonDB{
		rw: rw,
	}
}

func (db *JsonDB) Save(comics map[int]*comic.Comic) {
	encoder := json.NewEncoder(db.rw)
	err := encoder.Encode(comics)

	if err != nil {
		log.Println(fmt.Errorf("error while saving comics: %w", err))
	}
}
func (db *JsonDB) Read() map[int]*comic.Comic {
	comics := make(map[int]*comic.Comic)
	decoder := json.NewDecoder(db.rw)
	err := decoder.Decode(&comics)

	if err != nil {
		log.Println(fmt.Errorf("error while loading comics: %w", err))
	}

	return comics
}
