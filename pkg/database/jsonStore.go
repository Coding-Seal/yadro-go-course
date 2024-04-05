package database

import (
	"encoding/json"
	"io"
	"yadro-go-course/pkg/xkcd"
)

type DB interface {
	Save(map[int]*xkcd.Comic)
}

type JsonDB struct {
	w io.Writer
}

func NewJsonDB(w io.Writer) *JsonDB {
	return &JsonDB{
		w: w,
	}
}

func (db *JsonDB) Save(comics map[int]*xkcd.Comic) {
	encoder := json.NewEncoder(db.w)
	_ = encoder.Encode(comics) //FIXME
}
