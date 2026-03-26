package ratelimit

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string]int
	limit    int
}

func (r *RateLimiter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests = make(map[string]int)
}

func (r *RateLimiter) Allow(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests[ip]++
	return r.requests[ip] <= r.limit
}

func NewRateLimiter(limit int) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]int),
		limit:    limit,
	}

	go func() {
		for {
			time.Sleep(time.Minute)
			rl.Reset()
		}
	}()
	return rl
}
