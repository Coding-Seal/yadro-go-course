package comic

import (
	"context"
	"errors"
	"os"
	"yadro-go-course/internal/core/models"
	"yadro-go-course/internal/core/ports"
	"yadro-go-course/pkg/jsonl"
)

var ErrNotFound = errors.New("comic not found")
var ErrInternal = errors.New("internal error")

type JsonRepo struct {
	file *os.File
}

var _ ports.ComicsRepo = (*JsonRepo)(nil)

func NewJsonDB(file *os.File) *JsonRepo {
	return &JsonRepo{
		file: file,
	}
}

func (db *JsonRepo) Comic(ctx context.Context, id int) (models.Comic, error) {
	sc := jsonl.NewScanner(db.file)
	var readComic models.Comic
	for sc.Scan() {

		if err := sc.Json(&readComic); err != nil {
			return readComic, errors.Join(ErrInternal, err)
		}
		if readComic.ID == id {
			return readComic, nil
		}
	}

	return readComic, ErrNotFound
}
func (db *JsonRepo) Store(ctx context.Context, comic models.Comic) error {
	wr := jsonl.NewWriter(db.file)

	if err := wr.WriteJson(comic); err != nil {
		return errors.Join(ErrInternal, err)
	}

	if err := wr.Flush(); err != nil {
		return errors.Join(ErrInternal, err)
	}
	return nil
}
