package ratelimit

import (
	"sync"
	"time"
)

// MemoryLimiter is a small fixed-window limiter for local development.
type MemoryLimiter struct {
	limit  int
	window time.Duration
	hits   map[string]bucket
	mu     sync.Mutex
}

type bucket struct {
	count   int
	expires time.Time
}

// NewMemoryLimiter creates a limiter suitable for lightweight gateway endpoints.
func NewMemoryLimiter(limit int, window time.Duration) *MemoryLimiter {
	return &MemoryLimiter{limit: limit, window: window, hits: make(map[string]bucket)}
}

// Allow records one hit and returns false when the key exceeds the configured limit.
func (l *MemoryLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	current := l.hits[key]
	if now.After(current.expires) {
		current = bucket{expires: now.Add(l.window)}
	}

	current.count++
	l.hits[key] = current
	return current.count <= l.limit
}
