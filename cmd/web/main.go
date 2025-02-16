package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"yadro-go-course/web"
	"yadro-go-course/web/rest"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	opts := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	slog.NewTextHandler(os.Stdout, &opts)
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &opts)))

	api := rest.NewClient("http://api:8080")
	srv := http.Server{
		Addr:    "0.0.0.0:8090",
		Handler: web.Routes(api),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()
	slog.Info("Server Started", slog.String("url", fmt.Sprintf("http://%s", srv.Addr)))
	<-ctx.Done()
	slog.Info("Server Stopped")

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown Failed", slog.Any("error", err))
		os.Exit(1)
	} else {
		slog.Info("Server Exited Properly")
	}
}
