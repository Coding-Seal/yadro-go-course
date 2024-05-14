package comic

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"yadro-go-course/internal/core/models"
	"yadro-go-course/internal/core/ports"
)

//go:embed sqlite-migrations
var fs embed.FS

type SqliteStore struct {
	db *sql.DB
}

func NewSqliteStore(db *sql.DB, migratePath string) *SqliteStore {
	return &SqliteStore{
		db: db,
	}
}

var _ ports.ComicsRepo = (*SqliteStore)(nil)

func (s SqliteStore) Comic(ctx context.Context, id int) (models.Comic, error) {
	row := s.db.QueryRowContext(ctx, "SELECT * FROM comics WHERE comic_id =$1 LIMIT (1)", id)
	comic := models.Comic{}
	err := row.Scan(&comic.ID, &comic.Title, &comic.Date, &comic.ImgURL,
		&comic.News, &comic.SafeTitle, &comic.Transcription, &comic.AltTranscription, &comic.Link)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Comic{}, errors.Join(ports.ErrNotFound, err)
		}
		return models.Comic{}, errors.Join(ports.ErrInternal, err)
	}
	return comic, nil
}

func (s SqliteStore) Store(ctx context.Context, comic models.Comic) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO comics (comic_id, title, date, img_url, news, safe_title, transcription, alt_transcription, link) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		comic.ID, comic.Title, comic.Date, comic.ImgURL,
		comic.News, comic.SafeTitle, comic.Transcription, comic.AltTranscription, comic.Link)
	if err != nil {
		return errors.Join(ports.ErrInternal, err)
	}
	return nil
}

func (s SqliteStore) ComicsAll(ctx context.Context) ([]models.Comic, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT * FROM comics WHERE comic_id =$1")
	if err != nil {
		return nil, errors.Join(ports.ErrInternal, err)
	}
	defer rows.Close() // it seems unnecessary
	var comics []models.Comic
	for rows.Next() {
		comic := models.Comic{}
		err := rows.Scan(&comic.ID, &comic.Title, &comic.Date,
			&comic.ImgURL, &comic.News, &comic.SafeTitle, &comic.Transcription, &comic.AltTranscription, &comic.Link)
		if err != nil {
			return nil, errors.Join(ports.ErrInternal, err)
		}
		comics = append(comics, comic)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Join(ports.ErrInternal, err)
	}
	return comics, nil
}
