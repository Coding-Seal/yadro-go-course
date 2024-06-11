package rest

import (
	"context"
	"testing"
	"time"

	"yadro-go-course/config"
)

func TestApi(t *testing.T) {
	Api(nil, nil, nil, &config.Config{
		Server: config.Server{
			RateLimit:        1,
			ConcurrencyLimit: 1,
			TokenMaxTime:     time.Minute,
			DeleteEvery:      time.Minute,
		},
	}, context.Background())
}
