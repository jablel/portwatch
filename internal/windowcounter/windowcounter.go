// Package windowcounter provides a sliding-window event counter keyed by
// an arbitrary string. It is useful for tracking how many times a port
// event has fired within a rolling time window.
package windowcounter

import (
	"sync"
	"time"
)

// entry records a single timestamped occurrence.
type entry struct {
	at time.Time
}

// Counter tracks event counts inside a sliding time window.
type Counter struct {
	mu     sync.Mutex
	window time.Duration
	buckets map[string][]entry
}

// New returns a Counter that counts events within the given window duration.
// A zero or negative window means every recorded event is always counted.
func New(window time.Duration) *Counter {
	return &Counter{
		window:  window,
		buckets: make(map[string][]entry),
	}
}

// Add records one occurrence for key at the current time and returns the
// total count of events for that key that fall within the window.
func (c *Counter) Add(key string) int {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buckets[key] = append(c.buckets[key], entry{at: now})
	return c.countLocked(key, now)
}

// Count returns the number of events for key that fall within the current
// window without recording a new event.
func (c *Counter) Count(key string) int {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.countLocked(key, now)
}

// Reset removes all recorded events for key.
func (c *Counter) Reset(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.buckets, key)
}

// countLocked prunes stale entries and returns the remaining count.
// Must be called with c.mu held.
func (c *Counter) countLocked(key string, now time.Time) int {
	if c.window <= 0 {
		return len(c.buckets[key])
	}
	cutoff := now.Add(-c.window)
	entries := c.buckets[key]
	i := 0
	for i < len(entries) && entries[i].at.Before(cutoff) {
		i++
	}
	c.buckets[key] = entries[i:]
	return len(c.buckets[key])
}
