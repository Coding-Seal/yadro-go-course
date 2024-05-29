package web

import (
	"context"
	"testing"
	"time"

	"yadro-go-course/config"
)

func TestServer_SetupRoutes(t *testing.T) {
	srv := NewServer()
	srv.SetupRoutes(&config.Config{
		Server: config.Server{
			RateLimit:        1,
			ConcurrencyLimit: 1,
			TokenMaxTime:     time.Minute,
			DeleteEvery:      time.Minute,
		},
	}, context.Background())
}
