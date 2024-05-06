package web

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	"yadro-go-course/internal/adapters/web/handlers"
	"yadro-go-course/internal/adapters/web/middleware"
	"yadro-go-course/internal/core/services"
	"yadro-go-course/pkg/words"
)

func Run(cfg *config.Config) {
	// Logger
	opts := slog.HandlerOptions{
		Level: nil,
	}

	switch cfg.Logger.Level {
	case "debug":
		opts.Level = slog.LevelDebug
	case "info":
		opts.Level = slog.LevelInfo
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	default:
		log.Fatalln("Invalid log level")
	}

	var logHand slog.Handler

	switch cfg.Logger.Type {
	case "json":
		logHand = slog.NewJSONHandler(os.Stdout, &opts)
	case "text":
		logHand = slog.NewTextHandler(os.Stdout, &opts)
	default:
		log.Fatalln("Invalid log type")
	}

	slog.SetDefault(slog.New(logHand))

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

	err = comicFetcher.Update(context.Background())

	if err != nil {
		slog.Error("Error updating comics", slog.Any("error", err))
	}

	slog.Info("Missing comics fetched")

	mux := http.NewServeMux()

	mux.Handle("GET /pics", handlers.WrapHandler(handlers.Search(searchService)))
	mux.Handle("POST /update", handlers.WrapHandler(handlers.Update(comicFetcher)))

	st := middleware.Stack(middleware.AddRequestID, middleware.Logging)
	srv := http.Server{
		Addr:              fmt.Sprintf("localhost:%d", cfg.Server.Port),
		Handler:           st(mux),
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()
	slog.Info("Server Started", slog.String("url", srv.Addr))

	<-done
	slog.Info("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown Failed", slog.Any("error", err))
	}

	slog.Info("Server Exited Properly")
}
