package ratelimiter

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}
type RateLimiterPerUser[C comparable] struct {
	limit   int
	mu      sync.Mutex
	clients map[C]*client
}

func NewRateLimiterPerUser[C comparable](limit int, deleteEvery time.Duration, ctx context.Context) *RateLimiterPerUser[C] {
	limiter := &RateLimiterPerUser[C]{
		limit:   limit,
		clients: make(map[C]*client),
	}
	ticker := time.NewTicker(deleteEvery)

	go func() {
		for {
			select {
			case <-ticker.C:
				limiter.mu.Lock()
				for id, client := range limiter.clients {
					if time.Since(client.lastSeen) > deleteEvery {
						delete(limiter.clients, id)
					}
				}
				limiter.mu.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()

	return limiter
}

func (r *RateLimiterPerUser[C]) Allow(id C) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if c, ok := r.clients[id]; ok {
		c.lastSeen = time.Now()
	} else {
		r.clients[id] = &client{rate.NewLimiter(rate.Limit(r.limit), r.limit), time.Now()}
	}

	return r.clients[id].limiter.Allow()
}
