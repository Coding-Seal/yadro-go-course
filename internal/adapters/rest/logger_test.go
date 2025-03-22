package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"yadro-go-course/config"
)

func TestSetupLogger(t *testing.T) {
	SetupLogger(&config.Config{
		Logger: config.Logger{Type: "text", Level: "debug"},
	})
	SetupLogger(&config.Config{
		Logger: config.Logger{Type: "text", Level: "info"},
	})
	SetupLogger(&config.Config{
		Logger: config.Logger{Type: "json", Level: "warn"},
	})
	SetupLogger(&config.Config{
		Logger: config.Logger{Type: "json", Level: "error"},
	})
	assert.Panics(t, func() {
		SetupLogger(&config.Config{
			Logger: config.Logger{Type: "default", Level: "default"},
		})
	})
}
