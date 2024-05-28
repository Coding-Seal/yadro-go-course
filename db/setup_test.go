package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"yadro-go-course/config"
)

const testDB = "../test/db/test.db"

func TestConnect(t *testing.T) {
	cfg := &config.Config{}
	cfg.DB.Url = testDB
	db, err := Connect(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestMigrateUp(t *testing.T) {
	cfg := &config.Config{}
	cfg.DB.Url = testDB
	db, err := Connect(cfg)
	assert.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, MigrateDown(db, context.Background())) })
	assert.NoError(t, MigrateUp(db, context.Background()))
}

func TestMigrateDown(t *testing.T) {
	cfg := &config.Config{}
	cfg.DB.Url = testDB
	db, err := Connect(cfg)
	assert.NoError(t, err)
	assert.NoError(t, MigrateDown(db, context.Background()))
}
