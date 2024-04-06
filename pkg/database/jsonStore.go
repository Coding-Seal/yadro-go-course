package database

import (
	"encoding/json"
	"io"
	"yadro-go-course/pkg/comic"
)

type JsonDB struct {
	w io.ReadWriter
}

func NewJsonDB(w io.ReadWriter) *JsonDB {
	return &JsonDB{
		w: w,
	}
}

func (db *JsonDB) Save(comics map[int]*comic.Comic) {
	encoder := json.NewEncoder(db.w)
	_ = encoder.Encode(comics) //FIXME
}
func (db *JsonDB) Read() map[int]*comic.Comic {
	//FIXME
	return nil
}
