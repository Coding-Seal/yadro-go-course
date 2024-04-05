package database

import (
	"fmt"
	"io"
	"yadro-go-course/pkg/xkcd"
)

type TextDB struct {
	w io.Writer
}

func NewTextDB(w io.Writer) *TextDB {
	return &TextDB{
		w: w,
	}
}

func (db *TextDB) Save(comics map[int]*xkcd.Comic) {
	for _, comic := range comics {
		_, err := fmt.Fprintf(db.w, "id=%d title=\"%s\" imgURL=\"%s\" keywords=%v\n",
			comic.ID, comic.Title, comic.ImgURL, comic.Keywords)
		if err != nil {
			return
		}
	}
}
