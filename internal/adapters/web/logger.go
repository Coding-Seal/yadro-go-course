package web

import (
	"log"
	"log/slog"
	"os"
	"yadro-go-course/config"
)

func SetupLogger(cfg *config.Config) {
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
}
