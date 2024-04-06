package database

import (
	"fmt"
	"io"
	"yadro-go-course/pkg/comic"
)

type TextDB struct {
	w io.Writer
}

func NewTextDB(w io.Writer) *TextDB {
	return &TextDB{
		w: w,
	}
}

func (db *TextDB) Save(comics map[int]*comic.Comic) {
	for _, comic := range comics {
		_, err := fmt.Fprintf(db.w, "id=%d imgURL=\"%s\" keywords=%v\n",
			comic.ID, comic.ImgURL, comic.Keywords)
		if err != nil {
			return
		}
	}
}
