package ratelimiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}
type PerUser[C comparable] struct {
	limit       int
	mu          sync.Mutex
	clients     map[C]*client
	deleteEvery time.Duration
}

func NewRateLimiterPerUser[C comparable](limit int, deleteEvery time.Duration) *PerUser[C] {
	return &PerUser[C]{
		limit:       limit,
		clients:     make(map[C]*client),
		deleteEvery: deleteEvery,
	}
}

func (r *PerUser[C]) CleanUp() {
	r.mu.Lock()
	for id, client := range r.clients {
		if time.Since(client.lastSeen) > r.deleteEvery {
			delete(r.clients, id)
		}
	}
	r.mu.Unlock()
}

func (r *PerUser[C]) Allow(id C) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if c, ok := r.clients[id]; ok {
		c.lastSeen = time.Now()
	} else {
		r.clients[id] = &client{rate.NewLimiter(rate.Limit(r.limit), r.limit), time.Now()}
	}

	return r.clients[id].limiter.Allow()
}
