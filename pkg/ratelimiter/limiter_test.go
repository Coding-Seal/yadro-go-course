package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPerUser_Allow(t *testing.T) {
	limiter := NewRateLimiterPerUser[int](1, 3)
	assert.True(t, limiter.Allow(5))
	assert.False(t, limiter.Allow(5))
}

func TestPerUser_CleanUp(t *testing.T) {
	limiter := NewRateLimiterPerUser[int](1, 3)
	limiter.Allow(5)
	rl := limiter.clients[5].limiter
	limiter.clients[5] = &client{limiter: rl, lastSeen: time.Now().Add(time.Duration(-5) * time.Second)}
	limiter.CleanUp()
	assert.NotContains(t, limiter.clients, 5)
}
