// Package portexpiry tracks how long a port has been continuously open
// and emits an alert when it exceeds a configured maximum age.
package portexpiry

import (
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// Entry records when a port was first seen.
type Entry struct {
	FirstSeen time.Time
	Port      scanner.Port
}

// Tracker monitors port open-duration and reports expired ports.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]Entry
	maxAge  time.Duration
	now     func() time.Time
}

// New creates a Tracker that flags ports open longer than maxAge.
// A zero or negative maxAge disables expiry (Expired always returns nil).
func New(maxAge time.Duration) *Tracker {
	return &Tracker{
		entries: make(map[string]Entry),
		maxAge:  maxAge,
		now:     time.Now,
	}
}

func portKey(p scanner.Port) string {
	return p.Proto + ":" + p.Addr
}

// Observe records the current set of open ports.
// Ports not present in current are removed from tracking.
func (t *Tracker) Observe(current []scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()

	seen := make(map[string]struct{}, len(current))
	for _, p := range current {
		k := portKey(p)
		seen[k] = struct{}{}
		if _, ok := t.entries[k]; !ok {
			t.entries[k] = Entry{FirstSeen: t.now(), Port: p}
		}
	}

	for k := range t.entries {
		if _, ok := seen[k]; !ok {
			delete(t.entries, k)
		}
	}
}

// Expired returns all ports whose open duration exceeds maxAge.
// Returns nil when maxAge is not positive.
func (t *Tracker) Expired() []Entry {
	if t.maxAge <= 0 {
		return nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	var out []Entry
	for _, e := range t.entries {
		if now.Sub(e.FirstSeen) > t.maxAge {
			out = append(out, e)
		}
	}
	return out
}

// Age returns how long the given port has been tracked.
// Returns 0 if the port is not currently tracked.
func (t *Tracker) Age(p scanner.Port) time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	if e, ok := t.entries[portKey(p)]; ok {
		return t.now().Sub(e.FirstSeen)
	}
	return 0
}
