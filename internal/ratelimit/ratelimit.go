// Package ratelimit provides per-key cooldown-based rate limiting for port
// change events. It prevents alert storms when a port flaps rapidly.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks the last allowed time for each key and suppresses subsequent
// calls that arrive within the configured cooldown window.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
	now      func() time.Time
}

// New returns a Limiter with the given cooldown duration.
// A zero or negative cooldown means every call is allowed.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the key has not been seen within the cooldown window.
// The first call for any key always returns true.
func (l *Limiter) Allow(key string) bool {
	if l.cooldown <= 0 {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	if t, ok := l.last[key]; ok && now.Sub(t) < l.cooldown {
		return false
	}
	l.last[key] = now
	return true
}

// Reset removes the cooldown record for a single key, allowing the next call
// through immediately.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}

// ResetAll clears all cooldown records.
func (l *Limiter) ResetAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[string]time.Time)
}

// Len returns the number of keys currently tracked.
func (l *Limiter) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.last)
}
