// Package ratelimit provides a small in-memory attempt limiter. The site runs
// as a single process with a hundred-odd guests, so a map behind a mutex is
// the right size of solution.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter allows a fixed number of attempts per key within a moving window.
type Limiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
	limit    int
	window   time.Duration
	// now is swappable so the tests do not have to sleep.
	now func() time.Time
}

// New returns a limiter allowing limit attempts per key per window.
func New(limit int, window time.Duration) *Limiter {
	return &Limiter{
		mu:       sync.Mutex{},
		attempts: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
		now:      time.Now,
	}
}

// WithClock replaces the clock, for tests.
func (l *Limiter) WithClock(now func() time.Time) *Limiter {
	l.now = now

	return l
}

// Allow records an attempt and reports whether it is within the limit.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	cutoff := now.Add(-l.window)

	kept := make([]time.Time, 0, len(l.attempts[key])+1)

	for _, at := range l.attempts[key] {
		if at.After(cutoff) {
			kept = append(kept, at)
		}
	}

	if len(kept) >= l.limit {
		l.attempts[key] = kept

		return false
	}

	l.attempts[key] = append(kept, now)

	return true
}

// Reset forgets a key's attempts, called once a guest gets in.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.attempts, key)
}
