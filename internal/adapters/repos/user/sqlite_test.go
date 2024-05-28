package user

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

func TestSqliteRepo_UserID(t *testing.T) {
	testUser := models.User{Login: "Bob"}
	ctx := context.Background()
	conn, err := db.Connect(cfg)
	assert.NoError(t, err)
	assert.NoError(t, db.MigrateUp(conn, ctx))
	t.Cleanup(func() { assert.NoError(t, db.MigrateDown(conn, ctx)) })

	store := NewSqliteRepo(conn)
	assert.NoError(t, store.AddUser(ctx, &testUser))
	u, err := store.UserID(ctx, testUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, testUser, u)

	_, err = store.UserID(ctx, 30)
	assert.ErrorIs(t, err, ports.ErrNotFound)
}

func TestSqliteRepo_UserLogin(t *testing.T) {
	testUser := models.User{Login: "Bob"}
	ctx := context.Background()
	conn, err := db.Connect(cfg)
	assert.NoError(t, err)
	assert.NoError(t, db.MigrateUp(conn, ctx))

	t.Cleanup(func() { assert.NoError(t, db.MigrateDown(conn, ctx)) })

	store := NewSqliteRepo(conn)
	assert.NoError(t, store.AddUser(ctx, &testUser))

	u, err := store.UserLogin(ctx, testUser.Login)
	assert.NoError(t, err)
	assert.Equal(t, testUser, u)

	_, err = store.UserLogin(ctx, "Alice")
	assert.ErrorIs(t, err, ports.ErrNotFound)
}

func TestSqliteRepo_RemoveUser(t *testing.T) {
	testUser := models.User{Login: "Bob"}
	ctx := context.Background()
	conn, err := db.Connect(cfg)
	assert.NoError(t, err)
	assert.NoError(t, db.MigrateUp(conn, ctx))
	t.Cleanup(func() { assert.NoError(t, db.MigrateDown(conn, ctx)) })

	store := NewSqliteRepo(conn)
	assert.NoError(t, store.AddUser(ctx, &testUser))

	u, err := store.UserID(ctx, testUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, testUser, u)

	assert.NoError(t, store.RemoveUser(ctx, testUser.ID))
	_, err = store.UserID(ctx, testUser.ID)
	assert.ErrorIs(t, err, ports.ErrNotFound)
}
