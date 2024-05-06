package comic

import (
	"context"
	"errors"
	"os"
	"sync"
	"yadro-go-course/internal/core/models"
	"yadro-go-course/internal/core/ports"
	"yadro-go-course/pkg/jsonl"
)

var ErrNotFound = errors.New("comic not found")
var ErrInternal = errors.New("internal error")

type JsonRepo struct {
	file *os.File
	m    map[int]models.Comic
	mu   sync.RWMutex
}

var _ ports.ComicsRepo = (*JsonRepo)(nil)

func NewJsonDB(file *os.File) *JsonRepo {
	r := &JsonRepo{
		file: file,
		m:    make(map[int]models.Comic),
	}
	sc := jsonl.NewScanner(file)

	for sc.Scan() {
		var comic models.Comic
		err := sc.Json(&comic)

		if err != nil {
			break
		}

		r.m[comic.ID] = comic
	}

	return r
}

func (db *JsonRepo) Comic(ctx context.Context, id int) (models.Comic, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if comic, ok := db.m[id]; ok {
		return comic, nil
	}

	return models.Comic{}, ErrNotFound
}
func (db *JsonRepo) Store(ctx context.Context, comic models.Comic) error {
	db.mu.RLock()
	if _, ok := db.m[comic.ID]; ok {
		db.mu.RUnlock()
		return nil
	}
	db.mu.RUnlock()

	db.mu.Lock()
	defer db.mu.Unlock()
	wr := jsonl.NewWriter(db.file)

	if err := wr.WriteJson(comic); err != nil {
		return errors.Join(ErrInternal, err)
	}

	if err := wr.Flush(); err != nil {
		return errors.Join(ErrInternal, err)
	}

	return nil
}
