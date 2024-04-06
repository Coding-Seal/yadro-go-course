package database

import (
	"encoding/json"
	"io"
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
	_ = encoder.Encode(comics) //FIXME
}
func (db *JsonDB) Read() map[int]*comic.Comic {
	//FIXME
	return nil
}
