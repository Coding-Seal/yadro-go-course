package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	sqlitedr "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
	"yadro-go-course/config"
)

//go:embed migrations
var migrations embed.FS

func Connect(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.DB.Url)
	return db, err
}

func MigrateUp(db *sql.DB, ctx context.Context) error {
	src, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("error creating migration src: %w", err)
	}

	dr, err := sqlitedr.WithInstance(db, &sqlitedr.Config{
		MigrationsTable: sqlitedr.DefaultMigrationsTable,
		DatabaseName:    "comics",
		NoTxWrap:        false,
	})
	if err != nil {
		return fmt.Errorf("error creating driver: %w", err)
	}

	mgr, err := migrate.NewWithInstance("iofs", src, "sqlite3", dr)
	if err != nil {
		return fmt.Errorf("error creating migrate: %w", err)
	}

	stopMgr := mgr.GracefulStop

	go func() {
		<-ctx.Done()
		stopMgr <- true
	}()

	err = mgr.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error runnung migrations : %w", err)
	}

	return nil
}
