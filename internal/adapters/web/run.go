package web

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yadro-go-course/config"
	"yadro-go-course/internal/adapters/repos/comic"
	"yadro-go-course/internal/adapters/repos/fetcher"
	"yadro-go-course/internal/adapters/repos/search"
	"yadro-go-course/internal/core/ports"
	"yadro-go-course/internal/core/services"
	"yadro-go-course/pkg/words"
)

func Run(cfg *config.Config) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Logger
	SetupLogger(cfg)

	// repos
	// DB
	var db ports.ComicsRepo
	switch cfg.DB.Type {
	case "json":
		slog.Info("Opening JsonDB", slog.String("url", cfg.DB.Url))
		file, err := os.OpenFile(cfg.Url, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

		if err != nil {
			slog.Error("Error opening file", slog.String("url", cfg.DB.Url))
			os.Exit(1)
		}

		defer file.Close()
		db = comic.NewJsonDB(file)
	case "sqlite":
		slog.Info("Opening SQLiteDB", slog.String("url", cfg.DB.Url))
		con, err := sql.Open("sqlite3", cfg.DB.Url)
		if err != nil {
			slog.Error("Error opening SQLiteDB", slog.String("url", cfg.DB.Url), slog.Any("error", err))
			os.Exit(1)
		}
		db = comic.NewSqliteStore(con)
	}
	// fetcher
	fet := fetcher.NewFetcher(cfg.SourceURL, cfg.Parallel)
	// Index
	stopWordsFile, err := os.OpenFile(cfg.StopWordsFile, os.O_RDONLY, 0666)
	if err != nil {
		slog.Error("Error opening file", slog.String("stopWordsFile", cfg.StopWordsFile))
		os.Exit(1)
	}
	defer stopWordsFile.Close()

	ind := search.NewIndex(words.NewStemmer(words.ParseStopWords(stopWordsFile)))

	// services
	searchService := services.NewSearch(ind, db)
	comicFetcher := services.NewFetcher(fet, db, ind)

	slog.Info("Fetching missing comics")

	if err := comicFetcher.Update(ctx); err != nil {
		slog.Error("Error updating comics", slog.Any("error", err))
	}

	slog.Info("Missing comics fetched")

	slog.Info("Building index")
	err = ind.Build(ctx, db)

	if err != nil {
		slog.Error("Error building index", slog.Any("error", err))
	}

	c := cron.New()
	_, err = c.AddFunc(cfg.UpdateSpec, func() {
		if err := comicFetcher.Update(ctx); err != nil {
			slog.Error("Error updating comics", slog.Any("error", err))
		}
	})

	if err != nil {
		slog.Error("Invalid cron spec", slog.Any("error", err))
	}

	srv := http.Server{
		Addr:    fmt.Sprintf("localhost:%d", cfg.Server.Port),
		Handler: Routes(comicFetcher, searchService),
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()
	slog.Info("Server Started", slog.String("url", fmt.Sprintf("http://%s", srv.Addr)))

	<-done
	slog.Info("Server Stopped")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown Failed", slog.Any("error", err))
	}

	slog.Info("Server Exited Properly")
}
