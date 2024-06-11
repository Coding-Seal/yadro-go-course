package rest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/robfig/cron/v3"

	"yadro-go-course/config"
	"yadro-go-course/db"
	"yadro-go-course/internal/adapters/repos/comic"
	"yadro-go-course/internal/adapters/repos/fetcher"
	"yadro-go-course/internal/adapters/repos/search"
	"yadro-go-course/internal/adapters/repos/user"
	"yadro-go-course/internal/core/services"
	"yadro-go-course/pkg/words"
)

type Server struct {
	Srv        http.Server
	cron       *cron.Cron
	UserSrv    *services.UserService
	ComicSrv   *services.Comic
	FetcherSrv *services.Fetcher
	SearchSrv  *services.Search
}

func NewServer() *Server {
	return &Server{Srv: http.Server{}}
}

func (s *Server) SetupServices(cfg *config.Config, ctx context.Context) error {
	slog.Info("Opening SQLiteDB", slog.String("url", cfg.DB.Url))

	conn, err := db.Connect(cfg)
	if err != nil {
		return fmt.Errorf("error opening SQLiteDB: %w", err)
	}

	slog.Info("Running migrations")

	if err := db.MigrateUp(conn, ctx); err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	comicRepo := comic.NewSqliteStore(conn)
	userRepo := user.NewSqliteRepo(conn)
	fetcherRepo := fetcher.NewFetcher(cfg.SourceURL, cfg.Parallel)

	stopWordsFile, err := os.OpenFile(cfg.StopWordsFile, os.O_RDONLY, 0o666)
	if err != nil {
		return fmt.Errorf("error opening stop words file: %w", err)
	}

	index := search.NewIndex(words.NewStemmer(words.ParseStopWords(stopWordsFile)))
	stopWordsFile.Close()

	// services
	s.SearchSrv = services.NewSearch(index, comicRepo)
	s.FetcherSrv = services.NewFetcher(fetcherRepo, comicRepo, index)
	s.UserSrv = services.NewUserService(userRepo)
	s.ComicSrv = services.NewComicService(comicRepo)

	slog.Info("Fetching missing comics")

	if err := s.FetcherSrv.Update(ctx); err != nil {
		return fmt.Errorf("error fetching missing comics: %w", err)
	}

	slog.Info("Building index")

	if err := index.Build(ctx, comicRepo); err != nil {
		return fmt.Errorf("error building index: %w", err)
	}

	return nil
}

func (s *Server) SetupServer(cfg *config.Config, ctx context.Context) error {
	s.cron = cron.New()

	_, err := s.cron.AddFunc(cfg.UpdateSpec, func() {
		if err := s.FetcherSrv.Update(ctx); err != nil {
			slog.Error("Error updating comics", slog.Any("error", err))
		}
	})
	if err != nil {
		return fmt.Errorf("error setting up cron: %w", err)
	}

	go func() {
		<-ctx.Done()
		s.cron.Stop()
	}()

	s.Srv.Addr = fmt.Sprintf("localhost:%d", cfg.Server.Port)

	return nil
}

func (s *Server) SetupRoutes(cfg *config.Config, ctx context.Context) {
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", Api(s.FetcherSrv, s.SearchSrv, s.UserSrv, cfg, ctx)))
	s.Srv.Handler = mux
}

func (s *Server) Start() error {
	err := s.Srv.ListenAndServe()
	s.cron.Start()

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer func() {
		cancel()
	}()

	return s.Srv.Shutdown(ctx)
}
