package web

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/bcrypt"

	"yadro-go-course/internal/adapters/repos/user"
	"yadro-go-course/internal/core/models"

	"github.com/golang-migrate/migrate/v4"
	sqlitedr "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"

	"yadro-go-course/config"
	"yadro-go-course/internal/adapters/repos/comic"
	"yadro-go-course/internal/adapters/repos/fetcher"
	"yadro-go-course/internal/adapters/repos/search"
	"yadro-go-course/internal/core/ports"
	"yadro-go-course/internal/core/services"
	"yadro-go-course/pkg/words"
)

//go:embed migrations
var migrations embed.FS

func Run(cfg *config.Config) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Logger
	SetupLogger(cfg)

	// repos
	// DB
	var comicDB ports.ComicsRepo

	slog.Info("Opening SQLiteDB", slog.String("url", cfg.DB.Url))

	con, err := sql.Open("sqlite3", cfg.DB.Url)
	if err != nil {
		slog.Error("Error opening SQLiteDB", slog.String("url", cfg.DB.Url), slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Running migrations")

	src, err := iofs.New(migrations, "migrations")
	if err != nil {
		slog.Error("Error creating migration src", slog.Any("error", err))
		os.Exit(1)
	}

	dr, err := sqlitedr.WithInstance(con, &sqlitedr.Config{
		MigrationsTable: sqlitedr.DefaultMigrationsTable,
		DatabaseName:    "comics",
		NoTxWrap:        false,
	})
	if err != nil {
		slog.Error("Error creating driver", slog.Any("error", err))
		os.Exit(1)
	}

	mgr, err := migrate.NewWithInstance("iofs", src, "sqlite3", dr)
	if err != nil {
		slog.Error("Error creating migrate", slog.Any("error", err))
		os.Exit(1)
	}

	stopMgr := mgr.GracefulStop

	go func() {
		<-ctx.Done()
		stopMgr <- true
	}()

	err = mgr.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("Error running migrations", slog.Any("error", err))
		os.Exit(1)
	}

	comicDB = comic.NewSqliteStore(con)
	userDB := user.NewSqliteRepo(con)

	// fetcher
	fet := fetcher.NewFetcher(cfg.SourceURL, cfg.Parallel)
	// Index
	stopWordsFile, err := os.OpenFile(cfg.StopWordsFile, os.O_RDONLY, 0o666)
	if err != nil {
		slog.Error("Error opening file", slog.String("stopWordsFile", cfg.StopWordsFile))
		os.Exit(1)
	}
	defer stopWordsFile.Close()

	ind := search.NewIndex(words.NewStemmer(words.ParseStopWords(stopWordsFile)))

	// services
	searchService := services.NewSearch(ind, comicDB)
	comicFetcher := services.NewFetcher(fet, comicDB, ind)
	userService := services.NewUserService(userDB)

	slog.Info("Fetching missing comics")

	if err := comicFetcher.Update(ctx); err != nil {
		slog.Error("Error updating comics", slog.Any("error", err))
	}

	slog.Info("Missing comics fetched")

	slog.Info("Building index")

	err = ind.Build(ctx, comicDB)
	if err != nil {
		slog.Error("Error building index", slog.Any("error", err))
	}

	slog.Info("Adding users")

	// adding users
	pswd, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	err = userService.AddUser(ctx, &models.User{Login: "admin", Password: pswd, IsAdmin: true})
	if err != nil {
		slog.Error("Error adding admin", slog.Any("error", err))
	}

	pswd, _ = bcrypt.GenerateFromPassword([]byte("alice"), bcrypt.DefaultCost)
	err = userService.AddUser(ctx, &models.User{Login: "alice", Password: pswd, IsAdmin: true})
	if err != nil {
		slog.Error("Error adding user", slog.Any("error", err))
	}

	pswd, _ = bcrypt.GenerateFromPassword([]byte("bob"), bcrypt.DefaultCost)
	err = userService.AddUser(ctx, &models.User{Login: "bob", Password: pswd, IsAdmin: true})
	if err != nil {
		slog.Error("Error adding user", slog.Any("error", err))
	}

	// server stuff
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
		Handler: Routes(comicFetcher, searchService, userService, cfg.ConcurrencyLimit, cfg.RateLimit, cfg.TokenMaxTime, cfg.DeleteEvery, ctx),
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
