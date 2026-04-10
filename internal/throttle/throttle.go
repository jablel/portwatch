// Package throttle provides rate-limiting for alert notifications,
// preventing alert storms when many ports change in a short window.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks per-key event times and suppresses events that occur
// more frequently than the configured minimum interval.
type Throttle struct {
	mu       sync.Mutex
	interval time.Duration
	last     map[string]time.Time
}

// New creates a Throttle that allows at most one event per key per interval.
// A zero or negative interval disables throttling (all events pass through).
func New(interval time.Duration) *Throttle {
	return &Throttle{
		interval: interval,
		last:     make(map[string]time.Time),
	}
}

// Allow reports whether the event identified by key should be allowed through.
// It returns true the first time a key is seen and then only after the
// configured interval has elapsed since the last allowed event.
func (t *Throttle) Allow(key string) bool {
	if t.interval <= 0 {
		return true
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if last, ok := t.last[key]; ok {
		if now.Sub(last) < t.interval {
			return false
		}
	}

	t.last[key] = now
	return true
}

// Reset clears the recorded time for key, so the next call to Allow will
// immediately return true regardless of when the previous event occurred.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// Len returns the number of keys currently tracked.
func (t *Throttle) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.last)
}
