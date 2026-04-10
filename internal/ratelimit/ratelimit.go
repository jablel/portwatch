// Package ratelimit provides per-key rate limiting for alert suppression.
// It prevents alert flooding by enforcing a minimum interval between
// successive alerts for the same port or event key.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks the last emission time per key and enforces a cooldown.
type Limiter struct {
	mu       sync.Mutex
	last     map[string]time.Time
	cooldown time.Duration
	now      func() time.Time
}

// New creates a Limiter with the given cooldown duration.
// A zero or negative cooldown disables rate limiting (all events pass).
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		last:     make(map[string]time.Time),
		cooldown: cooldown,
		now:      time.Now,
	}
}

// Allow returns true if the key has not been seen within the cooldown window.
// If allowed, the key's timestamp is updated to now.
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

// Reset clears the recorded timestamp for a key, allowing it immediately.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}

// ResetAll clears all recorded timestamps.
func (l *Limiter) ResetAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[string]time.Time)
}

// Remaining returns the duration until the key is allowed again.
// Returns zero if the key is currently allowed.
func (l *Limiter) Remaining(key string) time.Duration {
	if l.cooldown <= 0 {
		return 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	t, ok := l.last[key]
	if !ok {
		return 0
	}
	elapsed := l.now().Sub(t)
	if elapsed >= l.cooldown {
		return 0
	}
	return l.cooldown - elapsed
}
