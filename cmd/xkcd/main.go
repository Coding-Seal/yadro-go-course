package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/bcrypt"

	"yadro-go-course/config"
	"yadro-go-course/internal/adapters/web"
	"yadro-go-course/internal/core/models"
)

func main() {
	var configPath string

	var port int

	flag.StringVar(&configPath, "c", "config.yaml", "Path to config file")
	flag.IntVar(&port, "p", 0, "Port to listen on")
	flag.Parse()

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	if port != 0 {
		cfg.Port = port
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	web.SetupLogger(cfg)
	srv := web.NewServer()

	if err := srv.SetupServices(cfg, ctx); err != nil {
		slog.Error("Setup services failed", slog.Any("error", err))
		os.Exit(1)
	}

	if err := srv.SetupServer(cfg, ctx); err != nil {
		slog.Error("Setup services failed", slog.Any("error", err))
		os.Exit(1)
	}

	srv.SetupRoutes(cfg, ctx)

	addUsers(srv, ctx)

	go func() {
		if err := srv.Start(); err != nil {
			slog.Error("server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()
	slog.Info("Server Started", slog.String("url", fmt.Sprintf("http://%s", srv.Srv.Addr)))
	<-ctx.Done()
	slog.Info("Server Stopped")

	if err := srv.Stop(ctx); err != nil {
		slog.Error("Server Shutdown Failed", slog.Any("error", err))
		os.Exit(1)
	} else {
		slog.Info("Server Exited Properly")
	}
}

func addUsers(srv *web.Server, ctx context.Context) {
	pswd, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)

	err := srv.UserSrv.AddUser(ctx, &models.User{Login: "admin", Password: pswd, IsAdmin: true})
	if err != nil {
		slog.Error("Error adding admin", slog.Any("error", err))
	}

	pswd, _ = bcrypt.GenerateFromPassword([]byte("alice"), bcrypt.DefaultCost)

	err = srv.UserSrv.AddUser(ctx, &models.User{Login: "alice", Password: pswd, IsAdmin: true})
	if err != nil {
		slog.Error("Error adding user", slog.Any("error", err))
	}

	pswd, _ = bcrypt.GenerateFromPassword([]byte("bob"), bcrypt.DefaultCost)

	err = srv.UserSrv.AddUser(ctx, &models.User{Login: "bob", Password: pswd, IsAdmin: true})
	if err != nil {
		slog.Error("Error adding user", slog.Any("error", err))
	}
}
