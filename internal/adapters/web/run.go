package web

import (
	"context"
	"errors"
	"fmt"
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
	slog.Info("Opening DB", slog.String("url", cfg.DB.Url))
	file, err := os.OpenFile(cfg.Url, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		slog.Error("Error opening file", slog.String("url", cfg.DB.Url))
		os.Exit(1)
	}

	defer file.Close()
	db := comic.NewJsonDB(file)
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

	slog.Info("Building index")
	ind.MustBuild(db, fet)
	slog.Info("Index built")

	// services
	searchService := services.NewSearch(ind, db)
	comicFetcher := services.NewFetcher(fet, db, ind)

	slog.Info("Fetching missing comics")

	if err := comicFetcher.Update(ctx); err != nil {
		slog.Error("Error updating comics", slog.Any("error", err))
	}

	slog.Info("Missing comics fetched")

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
