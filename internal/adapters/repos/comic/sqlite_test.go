package comic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"yadro-go-course/config"
	"yadro-go-course/db"
	"yadro-go-course/internal/core/models"
	"yadro-go-course/internal/core/ports"
)

var cfg = &config.Config{DB: config.DB{Url: "../../../../test/db/test.db"}}

func TestSqliteStore_Comic(t *testing.T) {
	testComic := models.Comic{
		ID:    1,
		Title: "test",
	}
	ctx := context.Background()
	conn, err := db.Connect(cfg)
	assert.NoError(t, err)
	assert.NoError(t, db.MigrateUp(conn, ctx))
	t.Cleanup(func() { assert.NoError(t, db.MigrateDown(conn, ctx)) })

	store := NewSqliteStore(conn)
	_, err = store.Comic(ctx, 0)

	assert.ErrorIs(t, err, ports.ErrNotFound)

	assert.NoError(t, store.Store(ctx, testComic))
	comic, err := store.Comic(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, testComic, comic)

	assert.ErrorIs(t, store.Store(ctx, testComic), ports.ErrInternal)
}

func TestSqliteStore_Comics(t *testing.T) {
	var testComics []models.Comic

	ctx := context.Background()
	conn, err := db.Connect(cfg)
	assert.NoError(t, err)
	assert.NoError(t, db.MigrateUp(conn, ctx))
	t.Cleanup(func() { assert.NoError(t, db.MigrateDown(conn, ctx)) })

	store := NewSqliteStore(conn)

	for i := 1; i <= 5; i++ {
		comic := models.Comic{ID: i, Title: "test"}
		testComics = append(testComics, comic)
		assert.NoError(t, store.Store(ctx, comic))
	}

	comics, err := store.ComicsAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, testComics, comics)
}
