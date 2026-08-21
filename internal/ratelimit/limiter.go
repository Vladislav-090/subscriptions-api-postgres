package ratelimit

import (
	"sync"
	"time"
)

type entry struct {
	count       int
	windowStart time.Time
}

type Limiter struct {
	mu       sync.Mutex
	attempts map[string]entry
	limit    int
	window   time.Duration
}

func NewLimiter(limit int, window time.Duration) *Limiter {
	l := &Limiter{
		attempts: make(map[string]entry),
		limit:    limit,
		window:   window,
	}

	go l.cleanupLoop(window)

	return l
}

func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	e, ok := l.attempts[key]
	if !ok || time.Since(e.windowStart) >= l.window {
		l.attempts[key] = entry{
			count:       1,
			windowStart: time.Now(),
		}
		return true
	}

	if e.count < l.limit {
		e.count++
		l.attempts[key] = e
		return true
	}

	return false
}

func (l *Limiter) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		for key, e := range l.attempts {
			if time.Since(e.windowStart) >= l.window {
				delete(l.attempts, key)
			}
		}
		l.mu.Unlock()
	}
}